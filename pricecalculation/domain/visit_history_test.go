package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"implementing-ddd-in-go/pricecalculation/domain"
)

func TestNumberOfVisitsInMonthOfLastVisit_NoVisits(t *testing.T) {
	history := domain.NewVisitHistory()

	count := history.NumberOfVisitsInMonthOfLastVisit()

	assert.Equal(t, 0, count)
}

func TestNumberOfVisitsInMonthOfLastVisit_2VisitsInTheSameMonth(t *testing.T) {
	history := domain.NewVisitHistory()

	lastMonth := time.Now().AddDate(0, -1, 0)
	thisMonth := time.Now()

	history.Add(domain.NewVisit("person-1", lastMonth))
	history.Add(domain.NewVisit("person-1", thisMonth))
	history.Add(domain.NewVisit("person-1", thisMonth))

	count := history.NumberOfVisitsInMonthOfLastVisit()

	assert.Equal(t, 2, count)
}

func TestNumberOfVisitsInMonthOfLastVisit_2VisitsDifferentMonthCountAs1(t *testing.T) {
	history := domain.NewVisitHistory()

	lastMonth := time.Now().AddDate(0, -1, 0)
	thisMonth := time.Now()

	history.Add(domain.NewVisit("person-1", lastMonth))
	history.Add(domain.NewVisit("person-1", thisMonth))

	count := history.NumberOfVisitsInMonthOfLastVisit()

	assert.Equal(t, 1, count)
}
