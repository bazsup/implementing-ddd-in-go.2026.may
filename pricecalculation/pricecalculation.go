package pricecalculation

import "implementing-ddd-in-go/pricecalculation/domain"

type ForGettingVisitorByID func(id string) (domain.ExternalVisitor, error)

type FractionInput struct {
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
}

func NewPriceCalculator(getExternalVisitorByID ForGettingVisitorByID) PriceCalculator {
	return PriceCalculator{getExternalVisitorByID}
}

func (c PriceCalculator) CalculatePrice(personID, visitID string, fractions []FractionInput) (CalculatedPrice, error) {
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

	visit := NewVisit(personID, visitID, droppedFractions)

	var total domain.Price
	for _, df := range visit.droppedFractions {
		total = total.Add(df.CalculatePrice())
	}
	return CalculatedPrice{
		PersonID:      visit.personID,
		VisitID:       visit.visitID,
		PriceAmount:   total.Amount(),
		PriceCurrency: total.Currency(),
	}, nil
}
