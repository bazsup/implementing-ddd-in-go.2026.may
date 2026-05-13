package handler

import (
	"github.com/rs/zerolog"
)

type forReinitializingContext func()

type Handler struct {
	logger                zerolog.Logger
	reinitializingContext forReinitializingContext
}

func New(logger zerolog.Logger, reinitializingContext forReinitializingContext) *Handler {
	return &Handler{logger: logger, reinitializingContext: reinitializingContext}
}
