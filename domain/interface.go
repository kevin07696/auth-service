package domain

import "context"

type AuthMessageBusser interface {
	SendRegisterEvent(ctx context.Context, event RegisterEvent) StatusCode
}

type LoginData struct {
	UserID         string
	hashedPassword Password
}

type AuthReader interface {
	GetLoginByName(ctx context.Context, name Username) (LoginData, StatusCode)
	GetLoginByEmail(ctx context.Context, email Email) (LoginData, StatusCode)
}

type AuthRequest struct {
	Username string
	Email    string
	Password string
}

type AuthResponse struct {
	UserID string
}

type AuthServicer interface {
	Register(ctx context.Context, request AuthRequest) (AuthResponse, StatusCode)
	Login(ctx context.Context, request AuthRequest) (AuthResponse, StatusCode)
}
