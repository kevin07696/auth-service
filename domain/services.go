package domain

import (
	"context"
	"fmt"
	"log/slog"
)

type AuthService struct {
	hasher Hasher
	repo   Repositor
}

func NewLoginService(hasher Hasher, repo Repositor) AuthService {
	return AuthService{
		hasher: hasher,
		repo:   repo,
	}
}

func (a AuthService) Register(ctx context.Context, request CreateLoginRequest) (response LoginResponse, status StatusCode) {
	var isValid bool

	var username Username
	username, isValid = NewUsername(request.Username)
	if !isValid {
		slog.Error("Failed username validation")
		status = StatusBadRequest
		return
	}

	var emailComponents EmailComponents
	emailComponents, isValid = NewEmailComponents(request.Email)
	if !isValid {
		slog.Error("Failed email validation")
		status = StatusBadRequest
		return
	}

	var password Password
	password, isValid = NewPassword(request.Password)
	if !isValid {
		slog.Error("Failed password validation")
		status = StatusBadRequest
		return
	}
	password = password.HashPassword(a.hasher)

	login := Login{
		Username:       string(username),
		Email:          string(emailComponents.ToStandardString()),
		HashedPassword: string(password),
	}

	status = a.repo.CreateLogin(ctx, &login)
	if status > 0 {
		return
	}

	response.LoginID = fmt.Sprint(login.ID)
	return
}

func (a AuthService) Login(ctx context.Context, request LoginRequest) (response LoginResponse, status StatusCode) {
	var isValid bool

	var password Password
	password, isValid = NewPassword(request.Password)
	if !isValid {
		status = StatusBadRequest
		return
	}

	var login Login

	var emailComponents EmailComponents
	emailComponents, isValid = NewEmailComponents(request.UserInput)
	if isValid {
		login, status = a.repo.GetLoginByEmail(ctx, emailComponents.ToStandardString())
		if status > 0 {
			return
		}
	} else {
		var username Username
		username, isValid = NewUsername(request.UserInput)
		if isValid {
			login, status = a.repo.GetLoginByUsername(ctx, username)
			if status > 0 {
				return
			}
		}
	}

	isValid = a.hasher.VerifyPassword(Password(login.HashedPassword), password)
	if !isValid {
		status = StatusUnauthorized
		return
	}

	response.LoginID = fmt.Sprint(login.ID)
	return
}
