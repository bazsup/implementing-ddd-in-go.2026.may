package pricecalculation_test

import (
	"implementing-ddd-in-go/pricecalculation"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePrice_ForConstructionWaste(t *testing.T) {
	// arrange
	droppedFraction := pricecalculation.NewDroppedFraction("Construction waste", pricecalculation.NewWeightFromKG(10))

	// act
	price := droppedFraction.CalculatePrice()

	// assert
	assert.Equal(t, 1.5, price.Amount())
}

func TestCalculatePrice_ForGreenWaste(t *testing.T) {
	// arrange
	droppedFraction := pricecalculation.NewDroppedFraction("Green waste", pricecalculation.NewWeightFromKG(10))

	// act
	price := droppedFraction.CalculatePrice()

	// assert
	assert.Equal(t, 1.0, price.Amount())
}
