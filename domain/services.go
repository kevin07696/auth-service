package domain

import (
	"context"

	"github.com/google/uuid"
)

type AuthService struct {
	hasher Hasher
	bus    AuthMessageBusser
	reader AuthReader
}

func NewAuthService(hasher Hasher, bus AuthMessageBusser, reader AuthReader) AuthService {
	return AuthService{
		hasher: hasher,
		bus:    bus,
		reader: reader,
	}
}

func (a AuthService) Register(ctx context.Context, request AuthRequest) (AuthResponse, StatusCode) {
	var isValid bool

	var username Username
	username, isValid = NewUsername(request.Username)
	if !isValid {
		return AuthResponse{}, StatusBadRequest
	}

	var emailComponents EmailComponents
	emailComponents, isValid = NewEmailComponents(request.Email)
	if !isValid {
		return AuthResponse{}, StatusBadRequest
	}

	var password Password
	password, isValid = NewPassword(request.Password)
	if !isValid {
		return AuthResponse{}, StatusBadRequest
	}
	password = password.HashPassword(a.hasher)

	registerEvent := RegisterEvent{
		UserID:               uuid.NewString(),
		Name:                 username,
		Password:             password,
		FullEmailAddress:     emailComponents.ToFullString(),
		StandardEmailAddress: emailComponents.ToStandardString(),
	}

	status := a.bus.SendRegisterEvent(ctx, registerEvent)
	if status > 0 {
		return AuthResponse{}, status
	}

	return AuthResponse{UserID: registerEvent.UserID}, StatusOK
}

func (a AuthService) Login(ctx context.Context, request AuthRequest) (response AuthResponse, status StatusCode) {
	var isValid bool

	var password Password
	password, isValid = NewPassword(request.Password)
	if !isValid {
		return AuthResponse{}, StatusBadRequest
	}

	var loginData LoginData

	var username Username
	username, isValid = NewUsername(request.Username)
	if isValid {
		loginData, status = a.reader.GetLoginByName(ctx, username)
		if status > 0 {
			return AuthResponse{}, status
		}
	} else {
		var emailComponents EmailComponents
		emailComponents, isValid = NewEmailComponents(request.Email)
		if !isValid {
			return AuthResponse{}, StatusBadRequest
		}
		loginData, status = a.reader.GetLoginByEmail(ctx, emailComponents.ToStandardString())
		if status > 0 {
			return AuthResponse{}, status
		}
	}

	isValid = a.hasher.VerifyPassword(loginData.hashedPassword, password)
	if !isValid {
		return AuthResponse{}, StatusUnauthorized
	}

	return AuthResponse{UserID: loginData.UserID}, StatusOK
}
