package domain

import "errors"

var ErrUnknownFractionType = errors.New("fraction type not known")
var ErrUnknownCity = errors.New("city not known")

const ConstructionWaste = "Construction waste"
const GreenWaste = "Green waste"

var cityFractionPrices = map[string]map[string]Price{
	"Pineville": {
		GreenWaste:        NewPriceFromUSDcents(10),
		ConstructionWaste: NewPriceFromUSDcents(15),
	},
	"Oak City": {
		GreenWaste:        NewPriceFromUSDcents(8),
		ConstructionWaste: NewPriceFromUSDcents(19),
	},
}

type FractionType struct {
	name  string
	price Price
}

func NewFractionTypeFromString(fractionType, city string) (FractionType, error) {
	fractionPrices, ok := cityFractionPrices[city]
	if !ok {
		return FractionType{}, ErrUnknownCity
	}
	price, ok := fractionPrices[fractionType]
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
