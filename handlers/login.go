package handlers

import (
	"context"
	"log/slog"

	"github.com/kevin07696/login-service/domain"
	"github.com/kevin07696/login-service/protos"
)

type LoginHandler struct {
	protos.UnimplementedLoginServer

	*Handler
}

func NewLoginHandler(handler *Handler) *LoginHandler {
	return &LoginHandler{Handler: handler}
}

func (h LoginHandler) Register(ctx context.Context, request *protos.CreateLoginRequest) (*protos.LoginResponse, error) {
	registerRequest := domain.CreateLoginRequest{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}

	registerResponse, status := h.Service.Register(ctx, registerRequest)
	if status > 0 {
		slog.Error(h.ErrorHandler[status].Error())
		return nil, h.ErrorHandler[status]
	}

	slog.Info("Register Login Response", "login_id", registerResponse.LoginID)

	return &protos.LoginResponse{UserId: registerResponse.LoginID}, h.ErrorHandler[domain.StatusOK]
}

func (h LoginHandler) Login(ctx context.Context, request *protos.LoginRequest) (*protos.LoginResponse, error) {
	loginRequest := domain.LoginRequest{
		UserInput: request.UserInput,
		Password:  request.Password,
	}

	loginResponse, status := h.Service.Login(ctx, loginRequest)
	if status > 0 {
		return nil, h.ErrorHandler[status]
	}

	return &protos.LoginResponse{UserId: loginResponse.LoginID}, h.ErrorHandler[domain.StatusOK]
}
