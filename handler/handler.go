package handler

import (
	"implementing-ddd-in-go/pricecalculation"
	"implementing-ddd-in-go/pricecalculation/domain"

	"github.com/rs/zerolog"
)

type forReinitializingContext func()

type Handler struct {
	logger                zerolog.Logger
	reinitializingContext forReinitializingContext
	getVisitorByID        pricecalculation.ForGettingVisitorByID
	visitHistory          *domain.VisitHistory
}

func NewHandler(
	logger zerolog.Logger,
	reinitializingContext forReinitializingContext,
	getVisitorByID pricecalculation.ForGettingVisitorByID,
	visitHistory *domain.VisitHistory,
) *Handler {
	return &Handler{logger: logger, reinitializingContext: reinitializingContext, getVisitorByID: getVisitorByID, visitHistory: visitHistory}
}
