package domain

import (
	"context"
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

func (a AuthService) Register(ctx context.Context, request LoginRequest) (response LoginResponse, status StatusCode) {
	var isValid bool

	var username Username
	username, isValid = NewUsername(request.Username)
	if !isValid {
		status = StatusBadRequest
		return
	}

	var emailComponents EmailComponents
	emailComponents, isValid = NewEmailComponents(request.Email)
	if !isValid {
		status = StatusBadRequest
		return
	}

	var password Password
	password, isValid = NewPassword(request.Password)
	if !isValid {
		status = StatusBadRequest
		return
	}
	password = password.HashPassword(a.hasher)

	login := Login{
		Username:             username,
		StandardEmailAddress: emailComponents.ToStandardString(),
		HashedPassword:       password,
	}

	status = a.repo.CreateLogin(ctx, &login)
	if status > 0 {
		return
	}

	response.LoginID = login.ID
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

	var username Username
	username, isValid = NewUsername(request.Username)
	if isValid {
		login, status = a.repo.GetLoginByUsername(ctx, username)
		if status > 0 {
			return
		}
	} else {
		var emailComponents EmailComponents
		emailComponents, isValid = NewEmailComponents(request.Email)
		if !isValid {
			status = StatusBadRequest
			return
		}
		login, status = a.repo.GetLoginByEmail(ctx, emailComponents.ToStandardString())
		if status > 0 {
			return
		}
	}

	isValid = a.hasher.VerifyPassword(login.HashedPassword, password)
	if !isValid {
		status = StatusUnauthorized
		return
	}

	response.LoginID = login.ID
	return
}
