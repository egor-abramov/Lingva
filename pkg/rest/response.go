package rest

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response[T any] struct {
	Status string `json:"status"`
	Error  string `json:"message,omitempty"`
	Data   *T     `json:"data,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK[T any](data T) Response[T] {
	return Response[T]{
		Status: StatusOK,
		Data:   &data,
	}
}

func OKEmpty[T any]() Response[T] {
	return Response[T]{
		Status: StatusOK,
	}
}

func Error(msg string) Response[any] {
	return Response[any]{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response[any] {
	var errMessages []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field '%s' is required", err.Field()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field '%s' is not valid", err.Field()))
		}
	}
	return Response[any]{
		Status: StatusError,
		Error:  strings.Join(errMessages, "; "),
	}
}
