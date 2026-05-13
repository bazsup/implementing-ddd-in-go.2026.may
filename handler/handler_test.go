package handler_test

import (
	"implementing-ddd-in-go/pricecalculation/domain"
)

var dummyReinitContext = func() {}

var dummyGetVisitorByID = func(id string) (domain.ExternalVisitor, error) {
	v, _ := domain.NewExternalVisitor("private", id, "Pine Street 1", "Pineville")
	return v, nil
}

var dummyGetVisitHistoryByPersonID = func(id string) *domain.VisitHistory {
	return domain.NewVisitHistory(id)
}

var dummySaveVisitHistory = func(*domain.VisitHistory) {}
