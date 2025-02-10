package domain

import "context"

type AuthEvent struct {
	UserID           string
	Name             Username
	FullEmailAddress Email
}

type AuthMessageBusser interface {
	SendRegisterEvent(ctx context.Context, event AuthEvent) StatusCode
}

type LoginData struct {
	UserID         string
	hashedPassword Password
}

type UserCredentials struct {
	UserID               string
	Username             Username
	StandardEmailAddress Email
	HashedPassword       Password
}

type AuthReaderWriter interface {
	GetLoginByName(ctx context.Context, name Username) (LoginData, StatusCode)
	GetLoginByEmail(ctx context.Context, email Email) (LoginData, StatusCode)
	CreateLogin(ctx context.Context, request UserCredentials) StatusCode
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
