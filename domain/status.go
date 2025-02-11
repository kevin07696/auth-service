package domain

type StatusCode uint8

const (
	StatusOK           StatusCode = 0
	StatusBadRequest   StatusCode = 1
	StatusUnauthorized StatusCode = 2
	StatusNotFound     StatusCode = 3
	StatusDuplicateKey StatusCode = 4
	StatusInternal     StatusCode = 5
)
