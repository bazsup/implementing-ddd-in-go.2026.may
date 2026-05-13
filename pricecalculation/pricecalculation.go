package pricecalculation

import (
	"implementing-ddd-in-go/pricecalculation/domain"
)

type ForGettingVisitorByID func(id string) (domain.ExternalVisitor, error)

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
	getExternalVisitorByID ForGettingVisitorByID
	visitHistory           *domain.VisitHistory
}

func NewPriceCalculator(getExternalVisitorByID ForGettingVisitorByID, visitHistory *domain.VisitHistory) PriceCalculator {
	return PriceCalculator{getExternalVisitorByID, visitHistory}
}

func (c PriceCalculator) CalculatePrice(personID, visitID, date string, fractions []RawDroppedFraction) (CalculatedPrice, error) {
	visitor, err := c.getExternalVisitorByID(personID)
	if err != nil {
		return CalculatedPrice{}, err
	}

	var droppedFractions []domain.DroppedFraction
	for _, f := range fractions {
		ft, err := domain.NewFractionTypeFromString(f.Type, visitor.City())
		if err != nil {
			return CalculatedPrice{}, err
		}
		droppedFractions = append(droppedFractions, domain.NewDroppedFraction(ft, domain.NewWeightFromKG(f.AmountKG)))
	}

	var total domain.Price
	for _, df := range droppedFractions {
		total = total.Add(df.CalculatePrice())
	}

	visit, err := domain.NewVisit(personID, date)
	if err != nil {
		return CalculatedPrice{}, err
	}

	c.visitHistory.Add(visit)
	if c.visitHistory.NumberOfVisitsInMonthOfLastVisit() >= 3 {
		total = total.AddFee(5)
	}

	return CalculatedPrice{
		PersonID:      personID,
		VisitID:       visitID,
		PriceAmount:   total.Amount(),
		PriceCurrency: total.Currency(),
	}, nil
}
