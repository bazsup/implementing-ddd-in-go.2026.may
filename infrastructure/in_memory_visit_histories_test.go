package infrastructure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"implementing-ddd-in-go/infrastructure"
	"implementing-ddd-in-go/pricecalculation/domain"
)

func TestInMemoryVisitHistories_GetByPersonID_CreatesNewHistory(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()

	vh := histories.GetByPersonID("person-1")

	assert.NotNil(t, vh)
	assert.Equal(t, "person-1", vh.PersonId())
}

func TestInMemoryVisitHistories_GetByPersonID_ReturnsFreshCopyEachTimeForSamePerson(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()

	vh1 := histories.GetByPersonID("person-1")
	vh2 := histories.GetByPersonID("person-1")

	assert.NotSame(t, vh1, vh2)
	assert.Equal(t, vh1.PersonId(), vh2.PersonId())
	assert.Equal(t, vh1.Version(), vh2.Version())
}

func TestInMemoryVisitHistories_GetByPersonID_ReturnsDifferentHistoriesForDifferentPersons(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()

	vh1 := histories.GetByPersonID("person-1")
	vh2 := histories.GetByPersonID("person-2")

	assert.NotSame(t, vh1, vh2)
}

func TestInMemoryVisitHistories_Save_PersistsHistory(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()
	vh := domain.NewVisitHistory("person-1")

	err := histories.Save(vh)

	assert.NoError(t, err)
	loaded := histories.GetByPersonID("person-1")
	assert.Equal(t, "person-1", loaded.PersonId())
	assert.Equal(t, vh.Version(), loaded.Version())
}

func TestInMemoryVisitHistories_Save_DetectsConcurrentModification(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()
	visitor, _ := domain.NewExternalVisitor("business", "person-1", "addr", "Pineville")
	ft, _ := domain.NewFractionTypeFromString(domain.GreenWaste)
	fractions := []domain.DroppedFraction{domain.NewDroppedFraction(ft, domain.NewWeightFromKG(10))}

	// Two requests each load an independent copy of the same history
	requestA := histories.GetByPersonID("person-1")
	requestB := histories.GetByPersonID("person-1")

	// Request A completes first: mutate and save (version 0 → 1)
	visit, _ := domain.NewVisit("2026-05-01", visitor)
	_, _ = requestA.CalculatePriceOfVisit(visit, fractions, domain.NewFeePolicy(visitor), domain.DefaultFractionPricingPolicy)
	assert.NoError(t, histories.Save(requestA))

	// Request B tries to save with stale version (still 0 → 1, but stored is already 1)
	visit2, _ := domain.NewVisit("2026-05-02", visitor)
	_, _ = requestB.CalculatePriceOfVisit(visit2, fractions, domain.NewFeePolicy(visitor), domain.DefaultFractionPricingPolicy)
	err := histories.Save(requestB)

	assert.EqualError(t, err, domain.ErrConcurrentModification.Error())
}

func TestInMemoryVisitHistories_Save_SucceedsOnSequentialSaves(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()
	visitor, _ := domain.NewExternalVisitor("business", "person-1", "addr", "Pineville")
	ft, _ := domain.NewFractionTypeFromString(domain.GreenWaste)
	fractions := []domain.DroppedFraction{domain.NewDroppedFraction(ft, domain.NewWeightFromKG(10))}

	// First request: load, mutate, save
	firstLoad := histories.GetByPersonID("person-1")
	visit, _ := domain.NewVisit("2026-05-01", visitor)
	_, _ = firstLoad.CalculatePriceOfVisit(visit, fractions, domain.NewFeePolicy(visitor), domain.DefaultFractionPricingPolicy)
	assert.NoError(t, histories.Save(firstLoad))

	// Second request: reload (gets version 1), mutate, save → should succeed
	secondLoad := histories.GetByPersonID("person-1")
	assert.Equal(t, 1, secondLoad.Version())
	visit2, _ := domain.NewVisit("2026-05-02", visitor)
	_, _ = secondLoad.CalculatePriceOfVisit(visit2, fractions, domain.NewFeePolicy(visitor), domain.DefaultFractionPricingPolicy)
	assert.NoError(t, histories.Save(secondLoad))

	assert.Equal(t, 2, histories.GetByPersonID("person-1").Version())
}

func TestInMemoryVisitHistories_Reset_ClearsAllHistories(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()
	original := histories.GetByPersonID("person-1")

	histories.Reset()

	fresh := histories.GetByPersonID("person-1")
	assert.NotSame(t, original, fresh)
}
