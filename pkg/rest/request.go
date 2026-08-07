package rest

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func Decode[T any](log *slog.Logger, w http.ResponseWriter, r *http.Request) (T, bool) {
	var req T

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, Error("empty request"))
			return req, false
		}
		log.Error("failed to read request body", "error", err.Error())
		render.JSON(w, r, Error("failed to read request body"))
		return req, false
	}

	if err := validator.New().Struct(req); err != nil {
		var validationErr validator.ValidationErrors
		errors.As(err, &validationErr)
		log.Error("invalid request", "error", validationErr.Error())
		render.JSON(w, r, ValidationError(validationErr))
		return req, false
	}
	return req, true
}
