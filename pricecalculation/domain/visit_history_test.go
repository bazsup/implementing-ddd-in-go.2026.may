package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "implementing-ddd-in-go/pricecalculation/domain"
)

func TestNumberOfVisitsInMonthOfLastVisit_NoVisits(t *testing.T) {
	history := NewVisitHistory("test-person")

	count := history.NumberOfVisitsInMonthOfLastVisit()

	assert.Equal(t, 0, count)
}

func TestNumberOfVisitsInMonthOfLastVisit_2VisitsInTheSameMonth(t *testing.T) {
	history := NewVisitHistory("test-person")
	visitor, _ := NewExternalVisitor(ExternalVisitorTypePrivate, "test-person", "addr", "Pineville")

	previousVisit, err := NewVisit("2026-04-01", visitor)
	assert.NoError(t, err, "failed to create visit in test setup")

	currentVisit, err := NewVisit("2026-05-01", visitor)
	assert.NoError(t, err, "failed to create visit in test setup")

	history.Add(previousVisit)
	history.Add(currentVisit)
	history.Add(currentVisit)

	count := history.NumberOfVisitsInMonthOfLastVisit()

	assert.Equal(t, 2, count)
}

func privateCustomerFeePolicy() FeePolicy {
	visitor, _ := NewExternalVisitor(ExternalVisitorTypePrivate, "test", "addr", "Pineville")
	return NewFeePolicy(visitor)
}

func businessCustomerFeePolicy() FeePolicy {
	visitor, _ := NewExternalVisitor("business", "test", "addr", "Pineville")
	return NewFeePolicy(visitor)
}

func TestCalculatePriceOfVisit_WhenItIsTheFirstVisitEver(t *testing.T) {
	history := NewVisitHistory("person-1")
	visitor, _ := NewExternalVisitor(ExternalVisitorTypePrivate, "person-1", "Pine Street 1", "Pineville")
	visit, _ := NewVisit("2026-05-14", visitor)
	ft, _ := NewFractionTypeFromString(GreenWaste)
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(10))}

	price, err := history.CalculatePriceOfVisit(visit, fractions, privateCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	assert.Equal(t, NewPriceFromUSDcents(100), price)
}

func TestCalculatePriceOfVisit_WhenItIsTheFirstVisitThisMonth(t *testing.T) {
	history := NewVisitHistory("person-1")
	visitor, _ := NewExternalVisitor(ExternalVisitorTypePrivate, "person-1", "Pine Street 1", "Pineville")
	previousVisit1, _ := NewVisit("2026-04-01", visitor)
	previousVisit2, _ := NewVisit("2026-04-15", visitor)
	history.Add(previousVisit1)
	history.Add(previousVisit2)

	currentVisit, _ := NewVisit("2026-05-14", visitor)
	ft, _ := NewFractionTypeFromString(GreenWaste)
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(10))}

	price, err := history.CalculatePriceOfVisit(currentVisit, fractions, privateCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	assert.Equal(t, NewPriceFromUSDcents(100), price)
}

func TestCalculatePriceOfVisit_WhenItIsTheThirdVisitThisMonth_ShouldHaveAdditionalFee5Percent(t *testing.T) {
	history := NewVisitHistory("person-1")
	visitor, _ := NewExternalVisitor(ExternalVisitorTypePrivate, "person-1", "Pine Street 1", "Pineville")
	visit1, _ := NewVisit("2026-05-01", visitor)
	visit2, _ := NewVisit("2026-05-07", visitor)
	history.Add(visit1)
	history.Add(visit2)

	thirdVisit, _ := NewVisit("2026-05-14", visitor)
	ft, _ := NewFractionTypeFromString(GreenWaste)
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(10))}

	price, err := history.CalculatePriceOfVisit(thirdVisit, fractions, privateCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	assert.Equal(t, NewPriceFromUSDcents(105), price)
}

func TestCalculatePriceOfVisit_BusinessCustomerHasNoAdditionalFeeOnThirdVisitThisMonth(t *testing.T) {
	history := NewVisitHistory("person-1")
	visitor, _ := NewExternalVisitor("business", "person-1", "Pine Street 1", "Pineville")
	visit1, _ := NewVisit("2026-05-01", visitor)
	visit2, _ := NewVisit("2026-05-07", visitor)
	history.Add(visit1)
	history.Add(visit2)

	thirdVisit, _ := NewVisit("2026-05-14", visitor)
	ft, _ := NewFractionTypeFromString(GreenWaste)
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(10))}

	price, err := history.CalculatePriceOfVisit(thirdVisit, fractions, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	assert.Equal(t, NewPriceFromUSDcents(120), price)
}

func TestCalculatePriceOfVisit_OakCityBusinessConstructionWaste_FirstVisit(t *testing.T) {
	history := NewVisitHistory("business-1")
	visitor, _ := NewExternalVisitor("business", "business-1", "Oak Avenue 1", "Oak City")
	visit, _ := NewVisit("2026-05-01", visitor)
	ft, _ := NewFractionTypeFromString(ConstructionWaste)
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(600))}

	price, err := history.CalculatePriceOfVisit(visit, fractions, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	// 600 * 21 = 12600 cents = $126.00
	assert.Equal(t, NewPriceFromUSDcents(12600), price)
}

func TestCalculatePriceOfVisit_OakCityBusinessConstructionWaste_SecondVisitCrossesThreshold(t *testing.T) {
	history := NewVisitHistory("business-1")
	visitor, _ := NewExternalVisitor("business", "business-1", "Oak Avenue 1", "Oak City")
	firstVisit, _ := NewVisit("2026-05-01", visitor)
	ft, _ := NewFractionTypeFromString(ConstructionWaste)
	firstFractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(600))}
	_, _ = history.CalculatePriceOfVisit(firstVisit, firstFractions, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	secondVisit, _ := NewVisit("2026-06-01", visitor)
	secondFractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(900))}

	price, err := history.CalculatePriceOfVisit(secondVisit, secondFractions, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	// 400 kg remaining in tier 1: 400 * 21 = 8400; 500 kg in tier 2: 500 * 29 = 14500; total = 22900 cents = $229.00
	assert.Equal(t, NewPriceFromUSDcents(22900), price)
}

func TestCalculatePriceOfVisit_OakCityBusinessConstructionWaste_ExemptionResetsNextYear(t *testing.T) {
	history := NewVisitHistory("business-1")
	visitor, _ := NewExternalVisitor("business", "business-1", "Oak Avenue 1", "Oak City")
	ft, _ := NewFractionTypeFromString(ConstructionWaste)

	// Drop 1100 kg in 2026 (100 kg in tier 2)
	firstVisit, _ := NewVisit("2026-05-01", visitor)
	_, _ = history.CalculatePriceOfVisit(firstVisit, []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(1100))}, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	// Drop 100 kg in 2027 — exemption resets, all in tier 1
	newYearVisit, _ := NewVisit("2027-01-01", visitor)
	price, err := history.CalculatePriceOfVisit(newYearVisit, []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(100))}, businessCustomerFeePolicy(), DefaultFractionPricingPolicy)

	assert.NoError(t, err)
	// 100 * 21 = 2100 cents = $21.00
	assert.Equal(t, NewPriceFromUSDcents(2100), price)
}

func TestNumberOfVisitsInMonthOfLastVisit_DifferentPersonNotCount(t *testing.T) {
	history := NewVisitHistory("test-person")
	visitor1, _ := NewExternalVisitor(ExternalVisitorTypePrivate, "person-1", "addr", "Pineville")
	visitor2, _ := NewExternalVisitor(ExternalVisitorTypePrivate, "person-2", "addr", "Pineville")

	visitID1, err := NewVisit("2026-05-01", visitor1)
	assert.NoError(t, err, "failed to create visit in test setup")

	visitID2, err := NewVisit("2026-05-01", visitor2)
	assert.NoError(t, err, "failed to create visit in test setup")

	history.Add(visitID1)
	history.Add(visitID2)
	history.Add(visitID1)

	count := history.NumberOfVisitsInMonthOfLastVisit()

	assert.Equal(t, 2, count)
}
