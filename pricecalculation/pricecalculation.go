package pricecalculation

import (
	"implementing-ddd-in-go/pricecalculation/domain"
)

type ForGettingVisitorByID func(id string) (domain.ExternalVisitor, error)
type ForSavingVisitHistories func(*domain.VisitHistory)
type ForGettingVisitHistoriesByPersonID func(string) *domain.VisitHistory
type ForResettingVisitHistories func()

type RawDroppedFraction struct {
	Type     string
	AmountKG uint
}

type CalculatedPrice struct {
	PersonID      string
	VisitID       string
	PriceAmount   float64
	PriceCurrency string
}

type PriceCalculator struct {
	getExternalVisitorByID    ForGettingVisitorByID
	getVisitHistoryByPersonID ForGettingVisitHistoriesByPersonID
	saveVisitHistory          ForSavingVisitHistories
	fractionPricingPolicy     domain.FractionPricingPolicy
}

func NewPriceCalculator(
	getExternalVisitorByID ForGettingVisitorByID,
	getVisitHistoryByPersonID ForGettingVisitHistoriesByPersonID,
	saveVisitHistory ForSavingVisitHistories,
	fractionPricingPolicy domain.FractionPricingPolicy,
) PriceCalculator {
	return PriceCalculator{getExternalVisitorByID, getVisitHistoryByPersonID, saveVisitHistory, fractionPricingPolicy}
}

func (c PriceCalculator) CalculatePrice(personID, visitID, date string, fractions []RawDroppedFraction) (CalculatedPrice, error) {
	visitor, err := c.getExternalVisitorByID(personID)
	if err != nil {
		return CalculatedPrice{}, err
	}

	var droppedFractions []domain.DroppedFraction
	for _, f := range fractions {
		ft, err := domain.NewFractionTypeFromString(f.Type)
		if err != nil {
			return CalculatedPrice{}, err
		}
		droppedFractions = append(droppedFractions, domain.NewDroppedFraction(ft, domain.NewWeightFromKG(f.AmountKG)))
	}

	visit, err := domain.NewVisit(date, visitor)
	if err != nil {
		return CalculatedPrice{}, err
	}

	visitHistory := c.getVisitHistoryByPersonID(personID)
	total, err := visitHistory.CalculatePriceOfVisit(visit, droppedFractions, domain.NewFeePolicy(visitor), c.fractionPricingPolicy)
	if err != nil {
		return CalculatedPrice{}, err
	}
	c.saveVisitHistory(visitHistory)

	return CalculatedPrice{
		PersonID:      personID,
		VisitID:       visitID,
		PriceAmount:   total.Amount(),
		PriceCurrency: total.Currency(),
	}, nil
}
