package pricecalculation

type DroppedFraction struct {
	AmountKg     uint
	FractionType string
}

type Visit struct {
	personID         string
	visitID          string
	droppedFractions []DroppedFraction
}

func NewVisit(personID, visitID string, droppedFractions []DroppedFraction) Visit {
	return Visit{personID: personID, visitID: visitID, droppedFractions: droppedFractions}
}
