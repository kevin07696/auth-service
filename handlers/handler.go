package handlers

import (
	"github.com/kevin07696/auth-service/domain"
	"github.com/kevin07696/auth-service/proto"
)

type Handler struct {
	proto.UnimplementedAuthServer

	Service      domain.AuthServicer
	ErrorHandler []error
}

func NewHandler(service domain.AuthServicer) *Handler {

	return &Handler{
		Service:      service,
		ErrorHandler: initErrorHandler(),
	}
}
