package errs

import (
	"net/http"
)

type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func NewUnauthorizedError(message string) error { // 401
	return AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func NewNotFoundError(message string) error { // 404
	return AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func NewConflictError(message string) error { // 409
	return AppError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

func NewUnexpectedError() error { // 500
	return AppError{
		Code:    http.StatusInternalServerError,
		Message: "unexpected error",
	}
}

func NewBadRequestError(message string) error { // 400
	return AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func NewValidationError(message string) error { // 422
	return AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}
