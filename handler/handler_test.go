package handler_test

import (
	"implementing-ddd-in-go/pricecalculation/domain"
)

var dummyReinitContext = func() {}

var dummyGetVisitorByID = func(id string) (domain.Customer, error) {
	addr, _ := domain.NewAddress("Pine Street 1", "Pineville")
	return domain.NewPrivateCustomer(id, addr), nil
}

var dummyGetVisitHistoryByPersonID = func(id string) *domain.VisitHistory {
	return domain.NewVisitHistory(id)
}

var dummySaveVisitHistory = func(*domain.VisitHistory) error { return nil }

var dummyFractionPricingPolicy = domain.DefaultFractionPricingPolicy

var dummyPublishPriceCalculated = func(domain.PriceCalculated) error { return nil }
