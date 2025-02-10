package handlers

import (
	"context"

	"github.com/kevin07696/auth-service/domain"
	"github.com/kevin07696/auth-service/proto"
)

type AuthHandler struct {
	
	*Handler 
}

func NewAuthHandler(handler *Handler) AuthHandler {
	return AuthHandler{Handler: handler}
}

func (a AuthHandler) Register(ctx context.Context, request *proto.AuthRequest) (*proto.AuthResponse, error) {
	registerRequest := domain.AuthRequest{
		Username: request.Username,
		Email: request.Email,
		Password: request.Password,
	}

	registerResponse, status := a.Service.Register(ctx, registerRequest)
	if status > 0 {
		return nil, a.ErrorHandler[status]
	}

	return &proto.AuthResponse{UserId: registerResponse.UserID}, a.ErrorHandler[domain.StatusOK]
}

func (a AuthHandler) Login(ctx context.Context, request *proto.AuthRequest) (*proto.AuthResponse, error) {
	loginRequest := domain.AuthRequest{
		Username: request.Username,
		Email: request.Email,
		Password: request.Password,
	}

	loginResponse, status := a.Service.Login(ctx, loginRequest)
	if status > 0 {
		return nil, a.ErrorHandler[status]
	}

	return &proto.AuthResponse{UserId: loginResponse.UserID}, a.ErrorHandler[domain.StatusOK]
}