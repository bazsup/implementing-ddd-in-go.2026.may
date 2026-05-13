package domain_test

import (
	"implementing-ddd-in-go/pricecalculation/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePrice(t *testing.T) {
	tests := []struct {
		fractionType   string
		city           string
		customerType   string
		weightKG       uint
		expectedAmount float64
	}{
		{domain.ConstructionWaste, "Pineville", "private", 10, 1.5},
		{domain.GreenWaste, "Pineville", "private", 10, 1.0},
		{domain.ConstructionWaste, "Oak City", "private", 10, 1.9},
		{domain.GreenWaste, "Oak City", "private", 10, 0.8},
		{domain.ConstructionWaste, "Pineville", "business", 10, 1.3},
		{domain.GreenWaste, "Pineville", "business", 10, 1.2},
		{domain.ConstructionWaste, "Oak City", "business", 10, 2.1},
		{domain.GreenWaste, "Oak City", "business", 10, 0.8},
	}

	for _, tt := range tests {
		t.Run(tt.fractionType+"_"+tt.city+"_"+tt.customerType, func(t *testing.T) {
			fractionType, _ := domain.NewFractionTypeFromString(tt.fractionType, tt.city, tt.customerType)
			weight := domain.NewWeightFromKG(tt.weightKG)
			droppedFraction := domain.NewDroppedFraction(fractionType, weight)

			price := droppedFraction.CalculatePrice()

			assert.Equal(t, tt.expectedAmount, price.Amount())
		})
	}
}

func TestNewFractionTypeFromString(t *testing.T) {
	tests := []struct {
		fractionType string
		city         string
		customerType string
		expectedErr  error
	}{
		{domain.ConstructionWaste, "Pineville", "private", nil},
		{domain.GreenWaste, "Pineville", "private", nil},
		{domain.ConstructionWaste, "Oak City", "private", nil},
		{domain.GreenWaste, "Oak City", "private", nil},
		{domain.ConstructionWaste, "Pineville", "business", nil},
		{domain.GreenWaste, "Pineville", "business", nil},
		{domain.ConstructionWaste, "Oak City", "business", nil},
		{domain.GreenWaste, "Oak City", "business", nil},
		{"Special waste", "Pineville", "private", domain.ErrUnknownFractionType},
		{"Special waste", "Oak City", "private", domain.ErrUnknownFractionType},
		{domain.GreenWaste, "Unknown City", "private", domain.ErrUnknownCity},
		{domain.GreenWaste, "Pineville", "unknown", domain.ErrUnknownCustomerType},
	}

	for _, tt := range tests {
		t.Run(tt.fractionType+"_"+tt.city+"_"+tt.customerType, func(t *testing.T) {
			_, err := domain.NewFractionTypeFromString(tt.fractionType, tt.city, tt.customerType)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
