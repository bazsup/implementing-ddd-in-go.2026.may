package pricecalculation

import "errors"

var InvalidWeight = errors.New("weight amount must not be negative")

type Weight struct {
	kilograms int
}

func NewWeightFromKG(kilograms int) (Weight, error) {
	if kilograms < 0 {
		return Weight{}, InvalidWeight
	}
	return Weight{kilograms: kilograms}, nil
}

func (w Weight) Amount() int {
	return w.kilograms
}

func (w Weight) Equals(other Weight) bool {
	return w.kilograms == other.kilograms
}
