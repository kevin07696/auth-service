package handlers

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"
)

type ServerOption func(*Server)

/*
For internal microservices/gRPC: Use 49152–65535 to avoid conflicts.
For local development: Anything above 1024 is safe (e.g., 3000, 8080, 50051 for gRPC).
*/
func WithPort(port int) ServerOption {
	if port < 1024 {
		slog.Error("Invalid: Safe ports need to be abobe 1024.", "port", port)
	}
	return func(s *Server) {
		s.port = port
	}
}

func WithNetwork(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

type Server struct {
	grpc    *grpc.Server
	port    int
	network string
}

func (s *Server) Server() *grpc.Server {
	return s.grpc
}

func (s *Server) Serve() {
	slog.Info("Starting gRPC server.", "port", s.port)
	addr := fmt.Sprintf(":%d", s.port)
	l, err := net.Listen(s.network, addr)
	if err != nil {
		slog.Error("Failed to create listener.", "error", err)
		os.Exit(1)
	}
	err = s.grpc.Serve(l)
	if err != nil {
		slog.Error("Failed to create server.", "error", err)
		os.Exit(1)
	}
	slog.Info("grpc server listening", "addr", addr)
}

func NewServer(opts ...ServerOption) *Server {
	grpcServer := grpc.NewServer()

	server := &Server{
		port:    8000,
		network: "tcp",
		grpc:    grpcServer,
	}

	// Apply server options
	for _, opt := range opts {
		opt(server)
	}

	return server
}
