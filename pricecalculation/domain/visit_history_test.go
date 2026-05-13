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
	// arrange
	history := NewVisitHistory("test-person")
	lastMonth := "2026-04-01"
	thisMonth := "2026-05-01"

	previousVisit, err := NewVisit("person-1", lastMonth)
	assert.NoError(t, err, "failed to create visit in test setup")

	currentVisit, err := NewVisit("person-1", thisMonth)
	assert.NoError(t, err, "failed to create visit in test setup")

	history.Add(previousVisit)
	history.Add(currentVisit)
	history.Add(currentVisit)

	// act
	count := history.NumberOfVisitsInMonthOfLastVisit()

	// assert
	assert.Equal(t, 2, count)
}

func TestCalculatePriceOfVisit_WhenItIsTheFirstVisitEver(t *testing.T) {
	history := NewVisitHistory("person-1")
	visit, _ := NewVisit("person-1", "2026-05-14")
	ft, _ := NewFractionTypeFromString("Green waste", "Pineville", "private")
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(10))}

	price := history.CalculatePriceOfVisit(visit, fractions)

	assert.Equal(t, NewPriceFromUSDcents(100), price)
}

func TestCalculatePriceOfVisit_WhenItIsTheFirstVisitThisMonth(t *testing.T) {
	history := NewVisitHistory("person-1")
	previousVisit1, _ := NewVisit("person-1", "2026-04-01")
	previousVisit2, _ := NewVisit("person-1", "2026-04-15")
	history.Add(previousVisit1)
	history.Add(previousVisit2)

	currentVisit, _ := NewVisit("person-1", "2026-05-14")
	ft, _ := NewFractionTypeFromString("Green waste", "Pineville", "private")
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(10))}

	price := history.CalculatePriceOfVisit(currentVisit, fractions)

	assert.Equal(t, NewPriceFromUSDcents(100), price)
}

func TestCalculatePriceOfVisit_WhenItIsTheThirdVisitThisMonth_ShouldHaveAdditionalFee5Percent(t *testing.T) {
	history := NewVisitHistory("person-1")
	visit1, _ := NewVisit("person-1", "2026-05-01")
	visit2, _ := NewVisit("person-1", "2026-05-07")
	history.Add(visit1)
	history.Add(visit2)

	thirdVisit, _ := NewVisit("person-1", "2026-05-14")
	ft, _ := NewFractionTypeFromString("Green waste", "Pineville", "private")
	fractions := []DroppedFraction{NewDroppedFraction(ft, NewWeightFromKG(10))}

	price := history.CalculatePriceOfVisit(thirdVisit, fractions)

	assert.Equal(t, NewPriceFromUSDcents(105), price)
}

func TestNumberOfVisitsInMonthOfLastVisit_DifferentPersonNotCount(t *testing.T) {
	history := NewVisitHistory("test-person")

	thisMonth := "2026-05-01"

	visitID1, err := NewVisit("person-1", thisMonth)
	assert.NoError(t, err, "failed to create visit in test setup")

	visitID2, err := NewVisit("person-2", thisMonth)
	assert.NoError(t, err, "failed to create visit in test setup")

	history.Add(visitID1)
	history.Add(visitID2)
	history.Add(visitID1)

	count := history.NumberOfVisitsInMonthOfLastVisit()

	assert.Equal(t, 2, count)
}
