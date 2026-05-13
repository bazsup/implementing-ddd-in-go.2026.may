package handler

import (
	"implementing-ddd-in-go/pricecalculation"

	"github.com/rs/zerolog"
)

type forReinitializingContext func()

type Handler struct {
	logger                    zerolog.Logger
	reinitializingContext     forReinitializingContext
	getVisitorByID            pricecalculation.ForGettingVisitorByID
	getVisitHistoryByPersonID pricecalculation.ForGettingVisitHistoriesByPersonID
	saveVisitHistory          pricecalculation.ForSavingVisitHistories
}

func NewHandler(
	logger zerolog.Logger,
	reinitializingContext forReinitializingContext,
	getVisitorByID pricecalculation.ForGettingVisitorByID,
	getVisitHistoryByPersonID pricecalculation.ForGettingVisitHistoriesByPersonID,
	saveVisitHistory pricecalculation.ForSavingVisitHistories,
) *Handler {
	return &Handler{
		logger:                    logger,
		reinitializingContext:     reinitializingContext,
		getVisitorByID:            getVisitorByID,
		getVisitHistoryByPersonID: getVisitHistoryByPersonID,
		saveVisitHistory:          saveVisitHistory,
	}
}
