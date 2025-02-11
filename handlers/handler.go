package handlers

type Handler struct {
	Service      LoginServicer
	ErrorHandler []error
}

func NewHandler(service LoginServicer) *Handler {

	return &Handler{
		Service:      service,
		ErrorHandler: initErrorHandler(),
	}
}
