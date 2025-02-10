package handlers

import (
	"github.com/kevin07696/auth-service/domain"
	"github.com/kevin07696/auth-service/proto"
)

type Handler struct {
	proto.UnimplementedAuthServer

	Service      domain.AuthServicer
	Reader       domain.AuthReader
	ErrorHandler []error
}

func NewHandler(service domain.AuthServicer, reader domain.AuthReader) *Handler {

	return &Handler{
		Service:      service,
		Reader:       reader,
		ErrorHandler: initErrorHandler(),
	}
}
