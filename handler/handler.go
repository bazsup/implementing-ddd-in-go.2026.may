package handler

import (
	"implementing-ddd-in-go/pricecalculation"
	"implementing-ddd-in-go/pricecalculation/domain"

	"github.com/rs/zerolog"
)

type forReinitializingContext func()

type Handler struct {
	logger                      zerolog.Logger
	reinitializingContext       forReinitializingContext
	getCustomerByPersonID       pricecalculation.ForGettingCustomerByPersonID
	getVisitHistoryByCustomerID pricecalculation.ForGettingVisitHistoriesByCustomerID
	saveVisitHistory            pricecalculation.ForSavingVisitHistories
	fractionPricingPolicy       domain.FractionPricingPolicy
	publishDomainEvent          pricecalculation.ForPublishingDomainEvents
}

func NewHandler(
	logger zerolog.Logger,
	reinitializingContext forReinitializingContext,
	getCustomerByPersonID pricecalculation.ForGettingCustomerByPersonID,
	getVisitHistoryByCustomerID pricecalculation.ForGettingVisitHistoriesByCustomerID,
	saveVisitHistory pricecalculation.ForSavingVisitHistories,
	fractionPricingPolicy domain.FractionPricingPolicy,
	publishDomainEvent pricecalculation.ForPublishingDomainEvents,
) *Handler {
	return &Handler{
		logger:                      logger,
		reinitializingContext:       reinitializingContext,
		getCustomerByPersonID:       getCustomerByPersonID,
		getVisitHistoryByCustomerID: getVisitHistoryByCustomerID,
		saveVisitHistory:            saveVisitHistory,
		fractionPricingPolicy:       fractionPricingPolicy,
		publishDomainEvent:          publishDomainEvent,
	}
}
