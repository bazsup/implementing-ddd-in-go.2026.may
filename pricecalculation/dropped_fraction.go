package pricecalculation

import "errors"

var ErrUnknownFractionType = errors.New("fraction type not known")

const ConstructionWaste = "Construction waste"
const GreenWaste = "Green waste"

var allowedFractionTypePrices = map[string]Price{
	"Green waste":        NewPriceFromUSDcents(10),
	"Construction waste": NewPriceFromUSDcents(15),
}

type FractionType struct {
	name  string
	price Price
}

func NewFractionTypeFromString(fractionType string) (FractionType, error) {
	if price, ok := allowedFractionTypePrices[fractionType]; ok {
		return FractionType{name: fractionType, price: price}, nil
	}
	return FractionType{}, ErrUnknownFractionType
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
