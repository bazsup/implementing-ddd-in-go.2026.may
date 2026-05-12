package handler

import "github.com/rs/zerolog"

type Handler struct {
	logger zerolog.Logger
}

func New(logger zerolog.Logger) *Handler {
	return &Handler{logger: logger}
}
