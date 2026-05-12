package pricecalculation

import "implementing-ddd-in-go/pricecalculation/domain"

type Visit struct {
	personID         string
	visitID          string
	droppedFractions []domain.DroppedFraction
}

func NewVisit(personID, visitID string, droppedFractions []domain.DroppedFraction) Visit {
	return Visit{personID: personID, visitID: visitID, droppedFractions: droppedFractions}
}
