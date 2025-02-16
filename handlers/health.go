package handlers

import (
	"context"
	"log"
	"time"

	health "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthHandler struct {
	health.UnimplementedHealthServer

	duration time.Duration
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check implements the Health Check RPC.
func (h *HealthHandler) Check(ctx context.Context, req *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	log.Println("Received health check request for service:", req.Service)
	return &health.HealthCheckResponse{
		Status: health.HealthCheckResponse_SERVING,
	}, nil
}

// Watch implements the streaming Health Check RPC.
func (h *HealthHandler) Watch(req *health.HealthCheckRequest, stream health.Health_WatchServer) error {
	log.Println("Started health watch stream for service:", req.Service)

	// Periodically send health status
	for {
		select {
		case <-stream.Context().Done():
			log.Println("Health watch stream closed by client")
			return nil
		default:
			err := stream.Send(&health.HealthCheckResponse{
				Status: health.HealthCheckResponse_SERVING,
			})
			if err != nil {
				log.Printf("Error sending health status update: %v", err)
				return err
			}
			time.Sleep(h.duration * time.Second)
		}
	}
}
