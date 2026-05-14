package infrastructure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"implementing-ddd-in-go/infrastructure"
	"implementing-ddd-in-go/pricecalculation/domain"
)

func mustBusinessCustomerInfra(id, street, city string) domain.Customer {
	addr, _ := domain.NewAddress(street, city)
	c, _ := domain.NewCustomer("business", id, addr, "")
	return c
}

func TestInMemoryVisitHistories_GetByCustomerID_CreatesNewHistory(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()

	vh := histories.GetByCustomerID("person-1")

	assert.NotNil(t, vh)
	assert.Equal(t, "person-1", vh.CustomerID())
}

func TestInMemoryVisitHistories_GetByCustomerID_ReturnsFreshCopyEachTime(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()

	vh1 := histories.GetByCustomerID("person-1")
	vh2 := histories.GetByCustomerID("person-1")

	assert.NotSame(t, vh1, vh2)
	assert.Equal(t, vh1.CustomerID(), vh2.CustomerID())
	assert.Equal(t, vh1.Version(), vh2.Version())
}

func TestInMemoryVisitHistories_GetByCustomerID_ReturnsDifferentHistoriesForDifferentCustomers(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()

	vh1 := histories.GetByCustomerID("person-1")
	vh2 := histories.GetByCustomerID("person-2")

	assert.NotSame(t, vh1, vh2)
}

func TestInMemoryVisitHistories_Save_PersistsHistory(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()
	vh := domain.NewVisitHistory("person-1")

	err := histories.Save(vh)

	assert.NoError(t, err)
	loaded := histories.GetByCustomerID("person-1")
	assert.Equal(t, "person-1", loaded.CustomerID())
	assert.Equal(t, vh.Version(), loaded.Version())
}

func TestInMemoryVisitHistories_Save_DetectsConcurrentModification(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()
	customer := mustBusinessCustomerInfra("person-1", "addr", "Pineville")
	ft, _ := domain.NewFractionTypeFromString(domain.GreenWaste)
	fractions := []domain.DroppedFraction{domain.NewDroppedFraction(ft, domain.NewWeightFromKG(10))}

	// Two requests each load an independent copy of the same history
	requestA := histories.GetByCustomerID(customer.ID())
	requestB := histories.GetByCustomerID(customer.ID())

	// Request A completes first: mutate and save (version 0 → 1)
	visit, _ := domain.NewVisit("2026-05-01", customer)
	_, _ = requestA.CalculatePriceOfVisit(visit, fractions, domain.NewFeePolicy(customer), domain.DefaultFractionPricingPolicy)
	assert.NoError(t, histories.Save(requestA))

	// Request B tries to save with stale version (still 0 → 1, but stored is already 1)
	visit2, _ := domain.NewVisit("2026-05-02", customer)
	_, _ = requestB.CalculatePriceOfVisit(visit2, fractions, domain.NewFeePolicy(customer), domain.DefaultFractionPricingPolicy)
	err := histories.Save(requestB)

	assert.EqualError(t, err, domain.ErrConcurrentModification.Error())
}

func TestInMemoryVisitHistories_Save_SucceedsOnSequentialSaves(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()
	customer := mustBusinessCustomerInfra("person-1", "addr", "Pineville")
	ft, _ := domain.NewFractionTypeFromString(domain.GreenWaste)
	fractions := []domain.DroppedFraction{domain.NewDroppedFraction(ft, domain.NewWeightFromKG(10))}

	// First request: load, mutate, save
	firstLoad := histories.GetByCustomerID(customer.ID())
	visit, _ := domain.NewVisit("2026-05-01", customer)
	_, _ = firstLoad.CalculatePriceOfVisit(visit, fractions, domain.NewFeePolicy(customer), domain.DefaultFractionPricingPolicy)
	assert.NoError(t, histories.Save(firstLoad))

	// Second request: reload (gets version 1), mutate, save → should succeed
	secondLoad := histories.GetByCustomerID(customer.ID())
	assert.Equal(t, 1, secondLoad.Version())
	visit2, _ := domain.NewVisit("2026-05-02", customer)
	_, _ = secondLoad.CalculatePriceOfVisit(visit2, fractions, domain.NewFeePolicy(customer), domain.DefaultFractionPricingPolicy)
	assert.NoError(t, histories.Save(secondLoad))

	assert.Equal(t, 2, histories.GetByCustomerID(customer.ID()).Version())
}

func TestInMemoryVisitHistories_Reset_ClearsAllHistories(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()
	original := histories.GetByCustomerID("person-1")

	histories.Reset()

	fresh := histories.GetByCustomerID("person-1")
	assert.NotSame(t, original, fresh)
}
