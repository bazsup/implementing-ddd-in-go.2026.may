package domain

type FractionPriceCalculator interface {
	CalculatePrice(droppedFraction DroppedFraction) Price
}

type FractionPriceCalculatorFactory func(alreadyDroppedThisYear Weight) FractionPriceCalculator

type FlatRatePriceCalculator struct {
	pricePerKG Price
}

func NewFlatRatePriceCalculator(pricePerKG Price) FlatRatePriceCalculator {
	return FlatRatePriceCalculator{pricePerKG: pricePerKG}
}

func NewFlatRatePriceCalculatorFactory(pricePerKG Price) FractionPriceCalculatorFactory {
	return func(_ Weight) FractionPriceCalculator {
		return NewFlatRatePriceCalculator(pricePerKG)
	}
}

func (c FlatRatePriceCalculator) CalculatePrice(df DroppedFraction) Price {
	return c.pricePerKG.Times(df.weight.Amount())
}

type TierBasedPriceCalculator struct {
	threshold       Weight
	tier1Calculator FractionPriceCalculator
	tier2Calculator FractionPriceCalculator
	alreadyUsed     Weight
}

func NewTierBasedPriceCalculator(threshold Weight, tier1, tier2 FractionPriceCalculator, alreadyUsed Weight) TierBasedPriceCalculator {
	return TierBasedPriceCalculator{threshold: threshold, tier1Calculator: tier1, tier2Calculator: tier2, alreadyUsed: alreadyUsed}
}

func NewTierBasedPriceCalculatorFactory(threshold Weight, tier1, tier2 FractionPriceCalculator) FractionPriceCalculatorFactory {
	return func(alreadyUsed Weight) FractionPriceCalculator {
		return NewTierBasedPriceCalculator(threshold, tier1, tier2, alreadyUsed)
	}
}

func (c TierBasedPriceCalculator) CalculatePrice(df DroppedFraction) Price {
	var remaining uint
	if c.threshold.Amount() > c.alreadyUsed.Amount() {
		remaining = c.threshold.Amount() - c.alreadyUsed.Amount()
	}

	total := df.weight.Amount()
	tier1Weight := min(total, remaining)
	tier2Weight := total - tier1Weight

	ft := df.fractionType
	tier1Price := c.tier1Calculator.CalculatePrice(NewDroppedFraction(ft, NewWeightFromKG(tier1Weight)))
	tier2Price := c.tier2Calculator.CalculatePrice(NewDroppedFraction(ft, NewWeightFromKG(tier2Weight)))

	return tier1Price.Add(tier2Price)
}
