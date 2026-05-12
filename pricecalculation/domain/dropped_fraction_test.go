package domain_test

import (
	"implementing-ddd-in-go/pricecalculation/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePrice_ForConstructionWaste(t *testing.T) {
	// arrange
	fractionType, _ := domain.NewFractionTypeFromString(domain.ConstructionWaste)
	weight := domain.NewWeightFromKG(10)
	droppedFraction := domain.NewDroppedFraction(fractionType, weight)

	// act
	price := droppedFraction.CalculatePrice()

	// assert
	assert.Equal(t, 1.5, price.Amount())
}

func TestCalculatePrice_ForGreenWaste(t *testing.T) {
	// arrange
	fractionType, _ := domain.NewFractionTypeFromString(domain.GreenWaste)
	weight := domain.NewWeightFromKG(10)
	droppedFraction := domain.NewDroppedFraction(fractionType, weight)

	// act
	price := droppedFraction.CalculatePrice()

	// assert
	assert.Equal(t, 1.0, price.Amount())
}

func TestConstructionWasteFromString(t *testing.T) {
	// act
	_, err := domain.NewFractionTypeFromString(domain.ConstructionWaste)

	// assert
	assert.NoError(t, err)
}

func TestGreenWasteFromString(t *testing.T) {
	// act
	_, err := domain.NewFractionTypeFromString(domain.GreenWaste)

	// assert
	assert.NoError(t, err)
}

func TestUnknownFractionTypeFromString(t *testing.T) {
	// act
	_, err := domain.NewFractionTypeFromString("Special wast")

	// assert
	assert.EqualError(t, err, domain.ErrUnknownFractionType.Error())
}
