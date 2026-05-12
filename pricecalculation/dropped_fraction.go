package pricecalculation

import "errors"

var ErrUnknownFractionType = errors.New("fraction type not known")

const ConstructionWaste = "Construction waste"
const GreenWaste = "Green waste"

var allowedFractionTypePrices = map[string]Price{
	"Green waste":        NewPriceFromUSD(10),
	"Construction waste": NewPriceFromUSD(15),
}

type FractionType struct {
	fractionType string
	price        Price
}

func NewFractionTypeFromString(fractionType string) (FractionType, error) {
	if price, ok := allowedFractionTypePrices[fractionType]; ok {
		return FractionType{fractionType: fractionType, price: price}, nil
	}
	return FractionType{}, ErrUnknownFractionType
}

type DroppedFraction struct {
	weight       Weight
	FractionType string
}

func NewDroppedFraction(fractionType string, weight Weight) DroppedFraction {
	return DroppedFraction{FractionType: fractionType, weight: weight}
}

func (df DroppedFraction) CalculatePrice() Price {
	return allowedFractionTypePrices[df.FractionType].Times(df.weight.Amount())
}
