package pricecalculation_test

import (
	"implementing-ddd-in-go/pricecalculation"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePrice_ForConstructionWaste(t *testing.T) {
	// arrange
	fractionType, _ := pricecalculation.NewFractionTypeFromString(pricecalculation.ConstructionWaste)
	weight := pricecalculation.NewWeightFromKG(10)
	droppedFraction := pricecalculation.NewDroppedFraction(fractionType, weight)

	// act
	price := droppedFraction.CalculatePrice()

	// assert
	assert.Equal(t, 1.5, price.Amount())
}

func TestCalculatePrice_ForGreenWaste(t *testing.T) {
	// arrange
	fractionType, _ := pricecalculation.NewFractionTypeFromString(pricecalculation.GreenWaste)
	weight := pricecalculation.NewWeightFromKG(10)
	droppedFraction := pricecalculation.NewDroppedFraction(fractionType, weight)

	// act
	price := droppedFraction.CalculatePrice()

	// assert
	assert.Equal(t, 1.0, price.Amount())
}

func TestConstructionWasteFromString(t *testing.T) {
	// act
	_, err := pricecalculation.NewFractionTypeFromString(pricecalculation.ConstructionWaste)

	// assert
	assert.NoError(t, err)
}

func TestGreenWasteFromString(t *testing.T) {
	// act
	_, err := pricecalculation.NewFractionTypeFromString(pricecalculation.GreenWaste)

	// assert
	assert.NoError(t, err)
}

func TestUnknownFractionTypeFromString(t *testing.T) {
	// act
	_, err := pricecalculation.NewFractionTypeFromString("Special wast")

	// assert
	assert.EqualError(t, err, pricecalculation.ErrUnknownFractionType.Error())
}
