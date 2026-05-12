package pricecalculation

type Visit struct {
	personID string
	visitID  string
}

func NewVisit(personID, visitID string) Visit {
	return Visit{personID: personID, visitID: visitID}
}
