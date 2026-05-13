package pricecalculation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"implementing-ddd-in-go/pricecalculation"
	"implementing-ddd-in-go/pricecalculation/domain"
)

func TestCalculatePrice_ForABusinessCustomerFromPineville(t *testing.T) {
	getVisitor := func(id string) (domain.ExternalVisitor, error) {
		return domain.NewExternalVisitor("business", id, "Pine Street 1", "Pineville")
	}
	calculator := pricecalculation.NewPriceCalculator(getVisitor, emptyVisitHistory, noopSaveVisitHistory)

	result, err := calculator.CalculatePrice("person-1", "visit-1", "2026-05-14", []pricecalculation.RawDroppedFraction{
		{Type: "Green waste", AmountKG: 10},
		{Type: "Construction waste", AmountKG: 5},
	})

	assert.NoError(t, err)
	assert.Equal(t, "person-1", result.PersonID)
	assert.Equal(t, "visit-1", result.VisitID)
	// Green waste: 10 kg * 10 cents = 100 cents; Construction waste: 5 kg * 15 cents = 75 cents; total = 175 cents = 1.75 USD
	assert.Equal(t, 1.75, result.PriceAmount)
	assert.Equal(t, "USD", result.PriceCurrency)
}

func TestCalculatePrice_ForAPrivateCustomerFromOakCity(t *testing.T) {
	getVisitor := func(id string) (domain.ExternalVisitor, error) {
		return domain.NewExternalVisitor("private", id, "Oak Avenue 2", "Oak City")
	}
	calculator := pricecalculation.NewPriceCalculator(getVisitor, emptyVisitHistory, noopSaveVisitHistory)

	result, err := calculator.CalculatePrice("person-2", "visit-2", "2026-05-14", []pricecalculation.RawDroppedFraction{
		{Type: "Green waste", AmountKG: 10},
		{Type: "Construction waste", AmountKG: 5},
	})

	assert.NoError(t, err)
	assert.Equal(t, "person-2", result.PersonID)
	assert.Equal(t, "visit-2", result.VisitID)
	// Green waste: 10 kg * 8 cents = 80 cents; Construction waste: 5 kg * 19 cents = 95 cents; total = 175 cents = 1.75 USD
	assert.Equal(t, 1.75, result.PriceAmount)
	assert.Equal(t, "USD", result.PriceCurrency)
}

func TestCalculatePrice_WithAdditionalFeeFor3VisitsInOneMonth(t *testing.T) {
	getVisitor := func(id string) (domain.ExternalVisitor, error) {
		return domain.NewExternalVisitor("private", id, "Pine Street 1", "Pineville")
	}
	getVisitHistory := func(id string) *domain.VisitHistory {
		history := domain.NewVisitHistory(id)
		visit1, _ := domain.NewVisit(id, "2026-05-01")
		visit2, _ := domain.NewVisit(id, "2026-05-07")
		history.Add(visit1)
		history.Add(visit2)
		return history
	}
	calculator := pricecalculation.NewPriceCalculator(getVisitor, getVisitHistory, noopSaveVisitHistory)

	result, err := calculator.CalculatePrice("person-1", "visit-3", "2026-05-14", []pricecalculation.RawDroppedFraction{
		{Type: "Green waste", AmountKG: 10},
	})

	assert.NoError(t, err)
	// Green waste: 10 kg * 10 cents = 100 cents; +5% fee = 105 cents = 1.05 USD
	assert.Equal(t, 1.05, result.PriceAmount)
}

var emptyVisitHistory = func(id string) *domain.VisitHistory {
	return domain.NewVisitHistory(id)
}

var noopSaveVisitHistory = func(*domain.VisitHistory) {}
