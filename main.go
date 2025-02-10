package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kevin07696/auth-service/adapters"
	"github.com/kevin07696/auth-service/domain"
	"github.com/kevin07696/auth-service/handlers"
	"github.com/kevin07696/auth-service/proto"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var hasher domain.Hasher
	var authReader domain.AuthReader
	var authMessageBusser domain.AuthMessageBusser
	var authServicer domain.AuthServicer
	var handler *handlers.Handler

	hasher = adapters.NewBcryptAdapter(14)

	authServicer = domain.NewAuthService(hasher, authMessageBusser, authReader)

	handler = handlers.NewHandler(authServicer, authReader)

	grpcServer, err := handlers.NewServer(handler)
	if err != nil {
		logger.Error("failed to init new grpc server", "error", err)
		os.Exit(1)
	}

	proto.RegisterAuthServer(grpcServer.Server(), handlers.NewAuthHandler(handler))

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

	os.Exit(0) // Exit gracefully
}
