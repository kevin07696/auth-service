package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kevin07696/login-service/adapters"
	"github.com/kevin07696/login-service/domain"
	"github.com/kevin07696/login-service/handlers"
	"github.com/kevin07696/login-service/proto"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var hash domain.Hasher
	var repository domain.Repositor
	var service handlers.LoginServicer
	var handler *handlers.Handler

	hash = adapters.NewBcryptAdapter(14)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		logger.Error("failed to init db server", "error", err)
		os.Exit(1)
	}
	domain.NewLoginRepository(db)

	service = domain.NewLoginService(hash, repository)

	handler = handlers.NewHandler(service)

	grpcServer, err := handlers.NewServer(handler)
	if err != nil {
		logger.Error("failed to init grpc server", "error", err)
		os.Exit(1)
	}

	health.RegisterHealthServer(grpcServer.Server(), handlers.NewHealthHandler())
	proto.RegisterLoginServer(grpcServer.Server(), handlers.NewLoginHandler(handler))

	reflection.Register(grpcServer.Server())

	go func() {
		err := grpcServer.Serve()
		if err != nil {
			logger.Error("error serving grpc server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for a termination signal
	sig := <-quit
	logger.Info("Received signal. Shutting down...", "signal", sig)

	grpcServer.Server().GracefulStop()
}
