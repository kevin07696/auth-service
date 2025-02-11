package domain

import (
	"context"
)

type Repositor interface {
	Migrate()
	CreateLogin(ctx context.Context, login *Login) (status StatusCode)
	GetLoginByUsername(ctx context.Context, username Username) (login Login, status StatusCode)
	GetLoginByEmail(ctx context.Context, emailAddress Email) (login Login, status StatusCode)
}

type LoginRequest struct {
	Username string
	Email    string
	Password string
}

type LoginResponse struct {
	LoginID string
}
