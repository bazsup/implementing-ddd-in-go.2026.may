package handler

import (
	"implementing-ddd-in-go/pricecalculation"

	"github.com/rs/zerolog"
)

type forReinitializingContext func()

type Handler struct {
	logger                zerolog.Logger
	reinitializingContext forReinitializingContext
	getVisitorByID        pricecalculation.ForGettingVisitorByID
}

func NewHandler(
	logger zerolog.Logger,
	reinitializingContext forReinitializingContext,
	getVisitorByID pricecalculation.ForGettingVisitorByID,
) *Handler {
	return &Handler{logger: logger, reinitializingContext: reinitializingContext, getVisitorByID: getVisitorByID}
}
