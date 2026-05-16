package service

import "fmt"

type AppError struct {
	HTTPStatus int
	Code       int
	Message    string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func BadRequest(message string) *AppError {
	return &AppError{HTTPStatus: 400, Code: 400, Message: message}
}

func NotFound(message string) *AppError {
	return &AppError{HTTPStatus: 404, Code: 404, Message: message}
}

func Conflict(message string) *AppError {
	return &AppError{HTTPStatus: 409, Code: 409, Message: message}
}

func Internal(message string, err error) *AppError {
	return &AppError{HTTPStatus: 500, Code: 500, Message: message, Err: err}
}
