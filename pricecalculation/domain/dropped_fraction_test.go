package domain_test

import (
	"implementing-ddd-in-go/pricecalculation/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePrice_ForConstructionWasteInPineville(t *testing.T) {
	// arrange
	fractionType, _ := domain.NewFractionTypeFromString(domain.ConstructionWaste, "Pineville")
	weight := domain.NewWeightFromKG(10)
	droppedFraction := domain.NewDroppedFraction(fractionType, weight)

	// act
	price := droppedFraction.CalculatePrice()

	// assert
	assert.Equal(t, 1.5, price.Amount())
}

func TestCalculatePrice_ForGreenWasteInPineville(t *testing.T) {
	// arrange
	fractionType, _ := domain.NewFractionTypeFromString(domain.GreenWaste, "Pineville")
	weight := domain.NewWeightFromKG(10)
	droppedFraction := domain.NewDroppedFraction(fractionType, weight)

	// act
	price := droppedFraction.CalculatePrice()

	// assert
	assert.Equal(t, 1.0, price.Amount())
}

func TestCalculatePrice_ForConstructionWasteInOakCity(t *testing.T) {
	// arrange
	fractionType, _ := domain.NewFractionTypeFromString(domain.ConstructionWaste, "Oak City")
	weight := domain.NewWeightFromKG(10)
	droppedFraction := domain.NewDroppedFraction(fractionType, weight)

	// act
	price := droppedFraction.CalculatePrice()

	// assert
	assert.Equal(t, 1.9, price.Amount())
}

func TestCalculatePrice_ForGreenWasteInOakCity(t *testing.T) {
	// arrange
	fractionType, _ := domain.NewFractionTypeFromString(domain.GreenWaste, "Oak City")
	weight := domain.NewWeightFromKG(10)
	droppedFraction := domain.NewDroppedFraction(fractionType, weight)

	// act
	price := droppedFraction.CalculatePrice()

	// assert
	assert.Equal(t, 0.8, price.Amount())
}

func TestNewFractionTypeFromString_ConstructionWasteInPineville(t *testing.T) {
	// act
	_, err := domain.NewFractionTypeFromString(domain.ConstructionWaste, "Pineville")

	// assert
	assert.NoError(t, err)
}

func TestNewFractionTypeFromString_GreenWasteInPineville(t *testing.T) {
	// act
	_, err := domain.NewFractionTypeFromString(domain.GreenWaste, "Pineville")

	// assert
	assert.NoError(t, err)
}

func TestNewFractionTypeFromString_ConstructionWasteInOakCity(t *testing.T) {
	// act
	_, err := domain.NewFractionTypeFromString(domain.ConstructionWaste, "Oak City")

	// assert
	assert.NoError(t, err)
}

func TestNewFractionTypeFromString_GreenWasteInOakCity(t *testing.T) {
	// act
	_, err := domain.NewFractionTypeFromString(domain.GreenWaste, "Oak City")

	// assert
	assert.NoError(t, err)
}

func TestNewFractionTypeFromString_UnknownFractionTypeInPineville(t *testing.T) {
	// act
	_, err := domain.NewFractionTypeFromString("Special waste", "Pineville")

	// assert
	assert.EqualError(t, err, domain.ErrUnknownFractionType.Error())
}

func TestNewFractionTypeFromString_UnknownFractionTypeInOakCity(t *testing.T) {
	// act
	_, err := domain.NewFractionTypeFromString("Special waste", "Oak City")

	// assert
	assert.EqualError(t, err, domain.ErrUnknownFractionType.Error())
}

func TestNewFractionTypeFromString_KnownFractionTypeInUnknownCity(t *testing.T) {
	// act
	_, err := domain.NewFractionTypeFromString("Green waste", "Unknown City")

	// assert
	assert.EqualError(t, err, domain.ErrUnknownCity.Error())
}

