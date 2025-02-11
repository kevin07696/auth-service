package handlers

import (
	"context"

	"github.com/kevin07696/login-service/domain"
	"github.com/kevin07696/login-service/proto"
)

type LoginHandler struct {
	proto.UnimplementedLoginServer

	*Handler
}

func NewLoginHandler(handler *Handler) *LoginHandler {
	return &LoginHandler{Handler: handler}
}

func (h LoginHandler) Register(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	registerRequest := domain.LoginRequest{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	registerResponse, status := h.Service.Register(ctx, registerRequest)
	if status > 0 {
		return nil, h.ErrorHandler[status]
	}

	return &proto.LoginResponse{UserId: registerResponse.LoginID}, h.ErrorHandler[domain.StatusOK]
}

func (h LoginHandler) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	loginRequest := domain.LoginRequest{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	loginResponse, status := h.Service.Login(ctx, loginRequest)
	if status > 0 {
		return nil, h.ErrorHandler[status]
	}

	return &proto.LoginResponse{UserId: loginResponse.LoginID}, h.ErrorHandler[domain.StatusOK]
}
