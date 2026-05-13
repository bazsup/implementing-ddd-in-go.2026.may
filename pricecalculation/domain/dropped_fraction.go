package domain

import "errors"

var ErrUnknownFractionType = errors.New("fraction type not known")

const ConstructionWaste = "Construction waste"
const GreenWaste = "Green waste"

type FractionType struct {
	name string
}

func NewFractionTypeFromString(fractionType string) (FractionType, error) {
	switch fractionType {
	case ConstructionWaste:
		return FractionType{name: ConstructionWaste}, nil
	case GreenWaste:
		return FractionType{name: GreenWaste}, nil
	}
	return FractionType{}, ErrUnknownFractionType
}

type DroppedFraction struct {
	fractionType FractionType
	weight       Weight
}

func NewDroppedFraction(fractionType FractionType, weight Weight) DroppedFraction {
	return DroppedFraction{fractionType: fractionType, weight: weight}
}
