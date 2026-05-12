package pricecalculation

type Weight struct {
	kilograms uint
}

func NewWeightFromKG(kilograms uint) Weight {
	return Weight{kilograms: kilograms}
}

func (w Weight) Amount() uint {
	return w.kilograms
}

func (w Weight) Equals(other Weight) bool {
	return w.kilograms == other.kilograms
}
