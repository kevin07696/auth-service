package handlers

import (
	"context"

	"github.com/kevin07696/login-service/domain"
)

type LoginServicer interface {
	Register(ctx context.Context, request domain.LoginRequest) (domain.LoginResponse, domain.StatusCode)
	Login(ctx context.Context, request domain.LoginRequest) (domain.LoginResponse, domain.StatusCode)
}
