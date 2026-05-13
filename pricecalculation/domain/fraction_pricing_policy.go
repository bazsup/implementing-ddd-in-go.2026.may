package domain

type PriceKey struct {
	City         string
	FractionType string
	VisitorType  string
}

type FractionPricingPolicy map[PriceKey]FractionPriceCalculatorFactory

func (p FractionPricingPolicy) CalculatorFor(fractionType, city, visitorType string, alreadyDropped Weight) (FractionPriceCalculator, error) {
	factory, ok := p[PriceKey{City: city, FractionType: fractionType, VisitorType: visitorType}]
	if !ok {
		return nil, ErrUnknownFractionType
	}
	return factory(alreadyDropped), nil
}

var DefaultFractionPricingPolicy = FractionPricingPolicy{
	{City: "Pineville", FractionType: GreenWaste, VisitorType: "private"}:         NewFlatRatePriceCalculatorFactory(NewPriceFromUSDcents(10)),
	{City: "Pineville", FractionType: ConstructionWaste, VisitorType: "private"}:  NewFlatRatePriceCalculatorFactory(NewPriceFromUSDcents(15)),
	{City: "Oak City", FractionType: GreenWaste, VisitorType: "private"}:          NewFlatRatePriceCalculatorFactory(NewPriceFromUSDcents(8)),
	{City: "Oak City", FractionType: ConstructionWaste, VisitorType: "private"}:   NewFlatRatePriceCalculatorFactory(NewPriceFromUSDcents(19)),
	{City: "Pineville", FractionType: GreenWaste, VisitorType: "business"}:        NewFlatRatePriceCalculatorFactory(NewPriceFromUSDcents(12)),
	{City: "Pineville", FractionType: ConstructionWaste, VisitorType: "business"}: NewFlatRatePriceCalculatorFactory(NewPriceFromUSDcents(13)),
	{City: "Oak City", FractionType: GreenWaste, VisitorType: "business"}:         NewFlatRatePriceCalculatorFactory(NewPriceFromUSDcents(8)),
	{City: "Oak City", FractionType: ConstructionWaste, VisitorType: "business"}: NewTierBasedPriceCalculatorFactory(
		NewWeightFromKG(1000),
		NewFlatRatePriceCalculator(NewPriceFromUSDcents(21)),
		NewFlatRatePriceCalculator(NewPriceFromUSDcents(29)),
	),
}
