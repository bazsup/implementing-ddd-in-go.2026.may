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

func TestInMemoryVisitHistories_GetByPersonID_ReturnsSameHistoryForSamePerson(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()

	vh1 := histories.GetByPersonID("person-1")
	vh2 := histories.GetByPersonID("person-1")

	assert.Same(t, vh1, vh2)
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

	histories.Save(vh)

	assert.Same(t, vh, histories.GetByPersonID("person-1"))
}

func TestInMemoryVisitHistories_Reset_ClearsAllHistories(t *testing.T) {
	histories := infrastructure.NewInMemoryVisitHistories()
	original := histories.GetByPersonID("person-1")

	histories.Reset()

	fresh := histories.GetByPersonID("person-1")
	assert.NotSame(t, original, fresh)
}
