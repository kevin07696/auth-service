package handlers

import (
	"github.com/kevin07696/auth-service/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func initErrorHandler() []error {
	handler := []error{}

	handler[domain.StatusOK] = status.Error(codes.OK, "Success!")
	handler[domain.StatusBadRequest] = status.Error(codes.InvalidArgument, "The login credentials are invalid.")
	handler[domain.StatusUnauthorized] = status.Error(codes.Unauthenticated, "The login credentials are invalid.")
	handler[domain.StatusInternalServer] = status.Error(codes.Internal, "Something went wrong. Please try again.")

	return handler
}
