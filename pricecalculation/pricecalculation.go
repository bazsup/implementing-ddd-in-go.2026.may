package pricecalculation

import (
	"implementing-ddd-in-go/pricecalculation/domain"
)

type ForGettingCustomerByPersonID func(personID string) (domain.Customer, error)
type ForSavingVisitHistories func(*domain.VisitHistory) error
type ForGettingVisitHistoriesByCustomerID func(customerID string) *domain.VisitHistory
type ForResettingVisitHistories func()
type ForPublishingDomainEvents func(domain.DomainEvent) error

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
	getCustomerByPersonID       ForGettingCustomerByPersonID
	getVisitHistoryByCustomerID ForGettingVisitHistoriesByCustomerID
	saveVisitHistory            ForSavingVisitHistories
	fractionPricingPolicy       domain.FractionPricingPolicy
	publishDomainEvent          ForPublishingDomainEvents
}

func NewPriceCalculator(
	getCustomerByPersonID ForGettingCustomerByPersonID,
	getVisitHistoryByCustomerID ForGettingVisitHistoriesByCustomerID,
	saveVisitHistory ForSavingVisitHistories,
	fractionPricingPolicy domain.FractionPricingPolicy,
	publishDomainEvent ForPublishingDomainEvents,
) PriceCalculator {
	return PriceCalculator{getCustomerByPersonID, getVisitHistoryByCustomerID, saveVisitHistory, fractionPricingPolicy, publishDomainEvent}
}

func (c PriceCalculator) CalculatePrice(personID, visitID, date string, fractions []RawDroppedFraction) (CalculatedPrice, error) {
	customer, err := c.getCustomerByPersonID(personID)
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

	visit, err := domain.NewVisit(date, customer)
	if err != nil {
		return CalculatedPrice{}, err
	}

	visitHistory := c.getVisitHistoryByCustomerID(customer.ID())
	total, err := visitHistory.CalculatePriceOfVisit(visit, droppedFractions, domain.NewFeePolicy(customer), c.fractionPricingPolicy)
	if err != nil {
		return CalculatedPrice{}, err
	}
	if err := c.saveVisitHistory(visitHistory); err != nil {
		return CalculatedPrice{}, err
	}

	var email string
	if bc, ok := customer.(domain.BusinessCustomer); ok {
		email = bc.Email()
	}

	if err := c.publishDomainEvent(domain.PriceCalculated{
		PersonID:      personID,
		VisitID:       visitID,
		CustomerType:  customer.Type(),
		PriceAmount:   total.Amount(),
		PriceCurrency: total.Currency(),
		Email:         email,
	}); err != nil {
		return CalculatedPrice{}, err
	}

	return CalculatedPrice{
		PersonID:      personID,
		VisitID:       visitID,
		PriceAmount:   total.Amount(),
		PriceCurrency: total.Currency(),
	}, nil
}
