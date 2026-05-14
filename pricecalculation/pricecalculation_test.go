package pricecalculation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"implementing-ddd-in-go/pricecalculation"
	"implementing-ddd-in-go/pricecalculation/domain"
)

func mustCustomer(visitorType, id, street, city string) domain.Customer {
	addr, _ := domain.NewAddress(street, city)
	c, _ := domain.NewCustomer(visitorType, id, addr, "")
	return c
}

func TestCalculatePrice_ForABusinessCustomerFromPineville(t *testing.T) {
	getCustomer := func(id string) (domain.Customer, error) {
		return mustCustomer("business", id, "Pine Street 1", "Pineville"), nil
	}
	calculator := pricecalculation.NewPriceCalculator(getCustomer, emptyVisitHistory, noopSaveVisitHistory, domain.DefaultFractionPricingPolicy)

	result, err := calculator.CalculatePrice("person-1", "visit-1", "2026-05-14", []pricecalculation.RawDroppedFraction{
		{Type: "Green waste", AmountKG: 10},
		{Type: "Construction waste", AmountKG: 5},
	})

	assert.NoError(t, err)
	assert.Equal(t, "person-1", result.PersonID)
	assert.Equal(t, "visit-1", result.VisitID)
	// Green waste: 10 kg * 12 cents = 120 cents; Construction waste: 5 kg * 13 cents = 65 cents; total = 185 cents = 1.85 USD
	assert.Equal(t, 1.85, result.PriceAmount)
	assert.Equal(t, "USD", result.PriceCurrency)
}

func TestCalculatePrice_ForAPrivateCustomerFromOakCity(t *testing.T) {
	getCustomer := func(id string) (domain.Customer, error) {
		return mustCustomer("private", id, "Oak Avenue 2", "Oak City"), nil
	}
	calculator := pricecalculation.NewPriceCalculator(getCustomer, emptyVisitHistory, noopSaveVisitHistory, domain.DefaultFractionPricingPolicy)

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
	getCustomer := func(id string) (domain.Customer, error) {
		return mustCustomer("private", id, "Pine Street 1", "Pineville"), nil
	}
	getVisitHistory := func(customerID string) *domain.VisitHistory {
		// For private customers, customerID == personID.
		customer := mustCustomer("private", customerID, "Pine Street 1", "Pineville")
		history := domain.NewVisitHistory(customerID)
		visit1, _ := domain.NewVisit("2026-05-01", customer)
		visit2, _ := domain.NewVisit("2026-05-07", customer)
		history.Add(visit1)
		history.Add(visit2)
		return history
	}
	calculator := pricecalculation.NewPriceCalculator(getCustomer, getVisitHistory, noopSaveVisitHistory, domain.DefaultFractionPricingPolicy)

	result, err := calculator.CalculatePrice("person-1", "visit-3", "2026-05-14", []pricecalculation.RawDroppedFraction{
		{Type: "Green waste", AmountKG: 10},
	})

	assert.NoError(t, err)
	// Green waste: 10 kg * 10 cents = 100 cents; +5% fee = 105 cents = 1.05 USD
	assert.Equal(t, 1.05, result.PriceAmount)
}

func TestCalculatePrice_BusinessCustomerHasNoAdditionalFeeFor3VisitsInOneMonth(t *testing.T) {
	getCustomer := func(id string) (domain.Customer, error) {
		return mustCustomer("business", id, "Pine Street 1", "Pineville"), nil
	}
	getVisitHistory := func(customerID string) *domain.VisitHistory {
		// Previous visits were made by different employees of the same business.
		// Their person IDs differ, but visitor.ID() is the same address-derived key.
		employeeA := mustCustomer("business", "employee-a", "Pine Street 1", "Pineville")
		history := domain.NewVisitHistory(customerID)
		visit1, _ := domain.NewVisit("2026-05-01", employeeA)
		visit2, _ := domain.NewVisit("2026-05-07", employeeA)
		history.Add(visit1)
		history.Add(visit2)
		return history
	}
	calculator := pricecalculation.NewPriceCalculator(getCustomer, getVisitHistory, noopSaveVisitHistory, domain.DefaultFractionPricingPolicy)

	result, err := calculator.CalculatePrice("person-1", "visit-3", "2026-05-14", []pricecalculation.RawDroppedFraction{
		{Type: "Green waste", AmountKG: 10},
	})

	assert.NoError(t, err)
	// Green waste: 10 kg * 12 cents = 120 cents; no additional fee for business customers
	assert.Equal(t, 1.20, result.PriceAmount)
}

func TestCalculatePrice_BusinessEmployeesShareYearlyWeightThreshold(t *testing.T) {
	// Employee A already dropped 600 kg of Construction waste this year.
	// Employee B now drops 600 kg — the business has 1200 kg total, so 200 kg must be priced at tier 2.
	employeeA := mustCustomer("business", "employee-a", "Oak Avenue 1", "Oak City")
	employeeB := mustCustomer("business", "employee-b", "Oak Avenue 1", "Oak City")

	getCustomer := func(id string) (domain.Customer, error) {
		return employeeB, nil
	}
	getVisitHistory := func(customerID string) *domain.VisitHistory {
		history := domain.NewVisitHistory(customerID)
		ft, _ := domain.NewFractionTypeFromString(domain.ConstructionWaste)
		visitA, _ := domain.NewVisit("2026-05-01", employeeA)
		_, _ = history.CalculatePriceOfVisit(visitA, []domain.DroppedFraction{domain.NewDroppedFraction(ft, domain.NewWeightFromKG(600))}, domain.NewFeePolicy(employeeA), domain.DefaultFractionPricingPolicy)
		return history
	}
	calculator := pricecalculation.NewPriceCalculator(getCustomer, getVisitHistory, noopSaveVisitHistory, domain.DefaultFractionPricingPolicy)

	result, err := calculator.CalculatePrice("employee-b", "visit-b", "2026-06-01", []pricecalculation.RawDroppedFraction{
		{Type: "Construction waste", AmountKG: 600},
	})

	assert.NoError(t, err)
	// 400 kg in tier 1: 400 * 21 = 8400; 200 kg in tier 2: 200 * 29 = 5800; total = 14200 cents = $142.00
	assert.Equal(t, 142.0, result.PriceAmount)
}

func TestCalculatePrice_ReturnsConcurrentModificationError(t *testing.T) {
	getCustomer := func(id string) (domain.Customer, error) {
		return mustCustomer("private", id, "Pine Street 1", "Pineville"), nil
	}
	staleHistory := func(customerID string) *domain.VisitHistory {
		return domain.NewVisitHistory(customerID)
	}
	saveAlwaysConflicts := func(*domain.VisitHistory) error {
		return domain.ErrConcurrentModification
	}
	calculator := pricecalculation.NewPriceCalculator(getCustomer, staleHistory, saveAlwaysConflicts, domain.DefaultFractionPricingPolicy)

	_, err := calculator.CalculatePrice("person-1", "visit-1", "2026-05-14", []pricecalculation.RawDroppedFraction{
		{Type: "Green waste", AmountKG: 10},
	})

	assert.ErrorIs(t, err, domain.ErrConcurrentModification)
}

var emptyVisitHistory = func(customerID string) *domain.VisitHistory {
	return domain.NewVisitHistory(customerID)
}

var noopSaveVisitHistory = func(*domain.VisitHistory) error { return nil }
