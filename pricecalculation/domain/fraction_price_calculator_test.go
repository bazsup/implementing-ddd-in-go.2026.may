package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "implementing-ddd-in-go/pricecalculation/domain"
)

func TestTierBasedPriceCalculator_AllWeightInTier1(t *testing.T) {
	ft, _ := NewFractionTypeFromString(ConstructionWaste)
	calc := NewTierBasedPriceCalculator(
		NewWeightFromKG(1000),
		NewFlatRatePriceCalculator(NewPriceFromUSDcents(21)),
		NewFlatRatePriceCalculator(NewPriceFromUSDcents(29)),
		NewWeightFromKG(0),
	)

	price := calc.CalculatePrice(NewDroppedFraction(ft, NewWeightFromKG(600)))

	// 600 * 21 = 12600 cents = $126.00
	assert.Equal(t, NewPriceFromUSDcents(12600), price)
}

func TestTierBasedPriceCalculator_WeightSpansBothTiers(t *testing.T) {
	ft, _ := NewFractionTypeFromString(ConstructionWaste)
	calc := NewTierBasedPriceCalculator(
		NewWeightFromKG(1000),
		NewFlatRatePriceCalculator(NewPriceFromUSDcents(21)),
		NewFlatRatePriceCalculator(NewPriceFromUSDcents(29)),
		NewWeightFromKG(600),
	)

	price := calc.CalculatePrice(NewDroppedFraction(ft, NewWeightFromKG(900)))

	// 400 * 21 + 500 * 29 = 8400 + 14500 = 22900 cents = $229.00
	assert.Equal(t, NewPriceFromUSDcents(22900), price)
}

func TestTierBasedPriceCalculator_AllWeightInTier2(t *testing.T) {
	ft, _ := NewFractionTypeFromString(ConstructionWaste)
	calc := NewTierBasedPriceCalculator(
		NewWeightFromKG(1000),
		NewFlatRatePriceCalculator(NewPriceFromUSDcents(21)),
		NewFlatRatePriceCalculator(NewPriceFromUSDcents(29)),
		NewWeightFromKG(1000),
	)

	price := calc.CalculatePrice(NewDroppedFraction(ft, NewWeightFromKG(500)))

	// 500 * 29 = 14500 cents = $145.00
	assert.Equal(t, NewPriceFromUSDcents(14500), price)
}
