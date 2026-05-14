package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "implementing-ddd-in-go/pricecalculation/domain"
)

func mustPrivateCustomer(id, street, city string) Customer {
	addr, _ := NewAddress(street, city)
	c, _ := NewCustomer("private", id, addr, "")
	return c
}

func mustBusinessCustomer(id, street, city string) Customer {
	addr, _ := NewAddress(street, city)
	c, _ := NewCustomer("business", id, addr, "")
	return c
}

func TestNumberOfVisitsInMonthOfLastVisit_NoVisits(t *testing.T) {
	history := NewVisitHistory("test-person")

	count := history.NumberOfVisitsInMonthOfLastVisit()

	assert.Equal(t, 0, count)
}

func TestNumberOfVisitsInMonthOfLastVisit_2VisitsInTheSameMonth(t *testing.T) {
	customer := mustPrivateCustomer("test-person", "addr", "Pineville")
	history := NewVisitHistory(customer.ID())

	previousVisit, err := NewVisit("2026-04-01", customer)
	assert.NoError(t, err, "failed to create visit in test setup")

	currentVisit, err := NewVisit("2026-05-01", customer)
	assert.NoError(t, err, "failed to create visit in test setup")

	history.Add(previousVisit)
	history.Add(currentVisit)
	history.Add(currentVisit)

	count := history.NumberOfVisitsInMonthOfLastVisit()

	assert.Equal(t, 2, count)
}

func privateCustomerFeePolicy() FeePolicy {
	customer := mustPrivateCustomer("test", "addr", "Pineville")
	return NewFeePolicy(customer)
}

func businessCustomerFeePolicy() FeePolicy {
	customer := mustBusinessCustomer("test", "addr", "Pineville")
	return NewFeePolicy(customer)
}

func TestCalculatePriceOfVisit_WhenItIsTheFirstVisitEver(t *testing.T) {
	customer := mustPrivateCustomer("person-1", "Pine Street 1", "Pineville")
	history := NewVisitHistory(customer.ID())
	visit, _ := NewVisit("2026-05-14", customer)
	ft, _ := NewFractionTypeFromString(GreenWaste)
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(10))}

	price, err := history.CalculatePriceOfVisit(visit, fractions, privateCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	assert.Equal(t, NewPriceFromUSDcents(100), price)
}

func TestCalculatePriceOfVisit_WhenItIsTheFirstVisitThisMonth(t *testing.T) {
	customer := mustPrivateCustomer("person-1", "Pine Street 1", "Pineville")
	history := NewVisitHistory(customer.ID())
	previousVisit1, _ := NewVisit("2026-04-01", customer)
	previousVisit2, _ := NewVisit("2026-04-15", customer)
	history.Add(previousVisit1)
	history.Add(previousVisit2)

	currentVisit, _ := NewVisit("2026-05-14", customer)
	ft, _ := NewFractionTypeFromString(GreenWaste)
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(10))}

	price, err := history.CalculatePriceOfVisit(currentVisit, fractions, privateCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	assert.Equal(t, NewPriceFromUSDcents(100), price)
}

func TestCalculatePriceOfVisit_WhenItIsTheThirdVisitThisMonth_ShouldHaveAdditionalFee5Percent(t *testing.T) {
	customer := mustPrivateCustomer("person-1", "Pine Street 1", "Pineville")
	history := NewVisitHistory(customer.ID())
	visit1, _ := NewVisit("2026-05-01", customer)
	visit2, _ := NewVisit("2026-05-07", customer)
	history.Add(visit1)
	history.Add(visit2)

	thirdVisit, _ := NewVisit("2026-05-14", customer)
	ft, _ := NewFractionTypeFromString(GreenWaste)
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(10))}

	price, err := history.CalculatePriceOfVisit(thirdVisit, fractions, privateCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	assert.Equal(t, NewPriceFromUSDcents(105), price)
}

func TestCalculatePriceOfVisit_BusinessCustomerHasNoAdditionalFeeOnThirdVisitThisMonth(t *testing.T) {
	customer := mustBusinessCustomer("person-1", "Pine Street 1", "Pineville")
	history := NewVisitHistory(customer.ID())
	visit1, _ := NewVisit("2026-05-01", customer)
	visit2, _ := NewVisit("2026-05-07", customer)
	history.Add(visit1)
	history.Add(visit2)

	thirdVisit, _ := NewVisit("2026-05-14", customer)
	ft, _ := NewFractionTypeFromString(GreenWaste)
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(10))}

	price, err := history.CalculatePriceOfVisit(thirdVisit, fractions, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	assert.Equal(t, NewPriceFromUSDcents(120), price)
}

func TestCalculatePriceOfVisit_OakCityBusinessConstructionWaste_FirstVisit(t *testing.T) {
	customer := mustBusinessCustomer("business-1", "Oak Avenue 1", "Oak City")
	history := NewVisitHistory(customer.ID())
	visit, _ := NewVisit("2026-05-01", customer)
	ft, _ := NewFractionTypeFromString(ConstructionWaste)
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(600))}

	price, err := history.CalculatePriceOfVisit(visit, fractions, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	// 600 * 21 = 12600 cents = $126.00
	assert.Equal(t, NewPriceFromUSDcents(12600), price)
}

func TestCalculatePriceOfVisit_OakCityBusinessConstructionWaste_SecondVisitCrossesThreshold(t *testing.T) {
	customer := mustBusinessCustomer("business-1", "Oak Avenue 1", "Oak City")
	history := NewVisitHistory(customer.ID())
	firstVisit, _ := NewVisit("2026-05-01", customer)
	ft, _ := NewFractionTypeFromString(ConstructionWaste)
	firstFractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(600))}
	_, _ = history.CalculatePriceOfVisit(firstVisit, firstFractions, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	secondVisit, _ := NewVisit("2026-06-01", customer)
	secondFractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(900))}

	price, err := history.CalculatePriceOfVisit(secondVisit, secondFractions, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	// 400 kg remaining in tier 1: 400 * 21 = 8400; 500 kg in tier 2: 500 * 29 = 14500; total = 22900 cents = $229.00
	assert.Equal(t, NewPriceFromUSDcents(22900), price)
}

func TestCalculatePriceOfVisit_OakCityBusinessConstructionWaste_ExemptionResetsNextYear(t *testing.T) {
	customer := mustBusinessCustomer("business-1", "Oak Avenue 1", "Oak City")
	history := NewVisitHistory(customer.ID())
	ft, _ := NewFractionTypeFromString(ConstructionWaste)

	// Drop 1100 kg in 2026 (100 kg in tier 2)
	firstVisit, _ := NewVisit("2026-05-01", customer)
	_, _ = history.CalculatePriceOfVisit(firstVisit, []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(1100))}, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	// Drop 100 kg in 2027 — exemption resets, all in tier 1
	newYearVisit, _ := NewVisit("2027-01-01", customer)
	price, err := history.CalculatePriceOfVisit(newYearVisit, []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(100))}, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	// 100 * 21 = 2100 cents = $21.00
	assert.Equal(t, NewPriceFromUSDcents(2100), price)
}

func TestNumberOfVisitsInMonthOfLastVisit_DifferentPrivateCustomerNotCount(t *testing.T) {
	customer1 := mustPrivateCustomer("person-1", "addr", "Pineville")
	customer2 := mustPrivateCustomer("person-2", "addr", "Pineville")

	// Two private customers sharing a history (unusual, but tests the ID-based grouping)
	history := NewVisitHistory(customer1.ID())

	visit1, err := NewVisit("2026-05-01", customer1)
	assert.NoError(t, err, "failed to create visit in test setup")

	visit2, err := NewVisit("2026-05-01", customer2)
	assert.NoError(t, err, "failed to create visit in test setup")

	history.Add(visit1)
	history.Add(visit2)
	history.Add(visit1)

	count := history.NumberOfVisitsInMonthOfLastVisit()

	assert.Equal(t, 2, count)
}

func TestCalculatePriceOfVisit_OakCityBusiness_TwoEmployeesShareThreshold(t *testing.T) {
	addr, _ := NewAddress("Oak Avenue 1", "Oak City")
	employeeA, _ := NewCustomer("business", "employee-a", addr, "")
	employeeB, _ := NewCustomer("business", "employee-b", addr, "")

	// Both employees share the same history (same business ID)
	history := NewVisitHistory(employeeA.ID())

	ft, _ := NewFractionTypeFromString(ConstructionWaste)

	// Employee A drops 600 kg
	visitA, _ := NewVisit("2026-05-01", employeeA)
	_, _ = history.CalculatePriceOfVisit(visitA, []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(600))}, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	// Employee B drops 600 kg — business has now dropped 1200 kg total, 200 kg should be in tier 2
	visitB, _ := NewVisit("2026-06-01", employeeB)
	price, err := history.CalculatePriceOfVisit(visitB, []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(600))}, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	// 400 kg in tier 1: 400 * 21 = 8400; 200 kg in tier 2: 200 * 29 = 5800; total = 14200 cents = $142.00
	assert.Equal(t, NewPriceFromUSDcents(14200), price)
}
