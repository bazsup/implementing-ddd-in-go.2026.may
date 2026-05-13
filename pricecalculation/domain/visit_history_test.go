package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "implementing-ddd-in-go/pricecalculation/domain"
)

func TestNumberOfVisitsInMonthOfLastVisit_NoVisits(t *testing.T) {
	history := NewVisitHistory()

	count := history.NumberOfVisitsInMonthOfLastVisit()

	assert.Equal(t, 0, count)
}

func TestNumberOfVisitsInMonthOfLastVisit_2VisitsInTheSameMonth(t *testing.T) {
	// arrange
	history := NewVisitHistory()
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

func TestNumberOfVisitsInMonthOfLastVisit_DifferentPersonNotCount(t *testing.T) {
	history := NewVisitHistory()

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
