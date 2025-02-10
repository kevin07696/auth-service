package domain

import (
	"context"

	"github.com/google/uuid"
)

type AuthService struct {
	hasher       Hasher
	bus          AuthMessageBusser
	readerWriter AuthReaderWriter
}

func NewAuthService(hasher Hasher, bus AuthMessageBusser, reader AuthReaderWriter) AuthService {
	return AuthService{
		hasher:       hasher,
		bus:          bus,
		readerWriter: reader,
	}
}

func (a AuthService) Register(ctx context.Context, request AuthRequest) (response AuthResponse, status StatusCode) {
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

	userCredentials := UserCredentials{
		UserID:               uuid.NewString(),
		Username:             username,
		StandardEmailAddress: emailComponents.ToStandardString(),
		HashedPassword:       password,
	}

	status = a.readerWriter.CreateLogin(ctx, userCredentials)
	if status > 0 {
		return
	}
	registerEvent := AuthEvent{
		UserID:           userCredentials.UserID,
		Name:             username,
		FullEmailAddress: emailComponents.ToFullString(),
	}

	status = a.bus.SendRegisterEvent(ctx, registerEvent)
	if status > 0 {
		return
	}

	response.UserID = userCredentials.UserID
	return
}

func (a AuthService) Login(ctx context.Context, request AuthRequest) (response AuthResponse, status StatusCode) {
	var isValid bool

	var password Password
	password, isValid = NewPassword(request.Password)
	if !isValid {
		status = StatusBadRequest
		return
	}

	var loginData LoginData

	var username Username
	username, isValid = NewUsername(request.Username)
	if isValid {
		loginData, status = a.readerWriter.GetLoginByName(ctx, username)
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
		loginData, status = a.readerWriter.GetLoginByEmail(ctx, emailComponents.ToStandardString())
		if status > 0 {
			return
		}
	}

	isValid = a.hasher.VerifyPassword(loginData.hashedPassword, password)
	if !isValid {
		status = StatusUnauthorized
		return
	}

	response.UserID = loginData.UserID
	return
}
