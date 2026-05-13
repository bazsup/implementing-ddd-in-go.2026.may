package handler_test

import "implementing-ddd-in-go/pricecalculation/domain"

var dummyReinitContext = func() {}

var dummyGetVisitorByID = func(id string) (domain.ExternalVisitor, error) {
	v, _ := domain.NewExternalVisitor("private", id, "Pine Street 1", "Pineville")
	return v, nil
}

var dummyVisitHistory = domain.NewVisitHistory()
