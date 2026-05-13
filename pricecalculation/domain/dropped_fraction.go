package domain

import "errors"

var ErrUnknownFractionType = errors.New("fraction type not known")
var ErrUnknownCity = errors.New("city not known")
var ErrUnknownCustomerType = errors.New("customer type not known")

const ConstructionWaste = "Construction waste"
const GreenWaste = "Green waste"

var fractionPrices = map[string]map[string]map[string]Price{
	"private": {
		"Pineville": {
			GreenWaste:        NewPriceFromUSDcents(10),
			ConstructionWaste: NewPriceFromUSDcents(15),
		},
		"Oak City": {
			GreenWaste:        NewPriceFromUSDcents(8),
			ConstructionWaste: NewPriceFromUSDcents(19),
		},
	},
	"business": {
		"Pineville": {
			GreenWaste:        NewPriceFromUSDcents(12),
			ConstructionWaste: NewPriceFromUSDcents(13),
		},
		"Oak City": {
			GreenWaste:        NewPriceFromUSDcents(8),
			ConstructionWaste: NewPriceFromUSDcents(21),
		},
	},
}

type FractionType struct {
	name  string
	price Price
}

func NewFractionTypeFromString(fractionType, city, customerType string) (FractionType, error) {
	cityPrices, ok := fractionPrices[customerType]
	if !ok {
		return FractionType{}, ErrUnknownCustomerType
	}
	wasteTypePrices, ok := cityPrices[city]
	if !ok {
		return FractionType{}, ErrUnknownCity
	}
	price, ok := wasteTypePrices[fractionType]
	if !ok {
		return FractionType{}, ErrUnknownFractionType
	}
	return FractionType{name: fractionType, price: price}, nil
}

type DroppedFraction struct {
	weight       Weight
	FractionType FractionType
}

func NewDroppedFraction(fractionType FractionType, weight Weight) DroppedFraction {
	return DroppedFraction{FractionType: fractionType, weight: weight}
}

func (df DroppedFraction) CalculatePrice() Price {
	return df.FractionType.price.Times(df.weight.Amount())
}
