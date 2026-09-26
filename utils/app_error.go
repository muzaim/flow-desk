package utils

import (
	"fmt"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Status    bool      `json:"status"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type AppError struct {
	StatusCode int
	Code       string
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewAppError(statusCode int, code string, message string, err error) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

func NewBadRequestError(code string, message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, code, message, err)
}

func NewNotFoundError(code string, message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, code, message, err)
}

func NewInternalServerError(code string, message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, code, message, err)
}
