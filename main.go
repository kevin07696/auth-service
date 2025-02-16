package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kevin07696/login-service/adapters"
	"github.com/kevin07696/login-service/domain"
	"github.com/kevin07696/login-service/handlers"
	"github.com/kevin07696/login-service/protos"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var hash domain.Hasher
	var service handlers.LoginServicer
	var handler *handlers.Handler

	hash = adapters.NewBcryptAdapter(14)

	db, err := gorm.Open(postgres.Open("postgres://root@cockroachdb:26257/defaultdb?sslmode=disable"), &gorm.Config{})
	if err != nil {
		logger.Error("failed to init db server", "error", err)
		os.Exit(1)
	}

	// Test the database connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get the underlying sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	logger.Info("Successfully connected to the database")

	var repository domain.Repositor = domain.NewLoginRepository(db)

	repository.Migrate()

	test := &domain.Login{
		Username:       "test-user",
		Email:          "test-user@example.com",
		HashedPassword: "124q35ws4tdrctvybu!",
	}

	if status := repository.CreateLogin(context.TODO(), test); status > 0 {
		logger.Error("Failed to insert sample login")
	}

	if _, status := repository.GetLoginByEmail(context.TODO(), domain.Email(test.Email)); status > 0 {
		logger.Error("Failed to get sample login")
	}

	service = domain.NewLoginService(hash, repository)

	handler = handlers.NewHandler(service)

	server := handlers.NewServer()

	health.RegisterHealthServer(server.Server(), handlers.NewHealthHandler())
	protos.RegisterLoginServer(server.Server(), handlers.NewLoginHandler(handler))

	reflection.Register(server.Server())

	go server.Serve()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for a termination signal
	sig := <-quit
	logger.Info("Received signal. Shutting down...", "signal", sig)

	server.Server().GracefulStop()
}
