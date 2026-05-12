package pricecalculation

type Visit struct {
	personID         string
	visitID          string
	droppedFractions []DroppedFraction
}

func NewVisit(personID, visitID string, droppedFractions []DroppedFraction) Visit {
	return Visit{personID: personID, visitID: visitID, droppedFractions: droppedFractions}
}
