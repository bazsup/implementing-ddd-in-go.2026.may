package domain

type PriceKey struct {
	City         string
	FractionType string
	VisitorType  string
}

type FractionPricingPolicy map[PriceKey]FractionPriceCalculator

func (p FractionPricingPolicy) CalculatorFor(fractionType, city, visitorType string) (FractionPriceCalculator, error) {
	calc, ok := p[PriceKey{City: city, FractionType: fractionType, VisitorType: visitorType}]
	if !ok {
		return nil, ErrUnknownFractionType
	}
	return calc, nil
}

var DefaultFractionPricingPolicy = FractionPricingPolicy{
	{City: "Pineville", FractionType: GreenWaste, VisitorType: "private"}:         NewFlatRatePriceCalculator(NewPriceFromUSDcents(10)),
	{City: "Pineville", FractionType: ConstructionWaste, VisitorType: "private"}:  NewFlatRatePriceCalculator(NewPriceFromUSDcents(15)),
	{City: "Oak City", FractionType: GreenWaste, VisitorType: "private"}:          NewFlatRatePriceCalculator(NewPriceFromUSDcents(8)),
	{City: "Oak City", FractionType: ConstructionWaste, VisitorType: "private"}:   NewFlatRatePriceCalculator(NewPriceFromUSDcents(19)),
	{City: "Pineville", FractionType: GreenWaste, VisitorType: "business"}:        NewFlatRatePriceCalculator(NewPriceFromUSDcents(12)),
	{City: "Pineville", FractionType: ConstructionWaste, VisitorType: "business"}: NewFlatRatePriceCalculator(NewPriceFromUSDcents(13)),
	{City: "Oak City", FractionType: GreenWaste, VisitorType: "business"}:         NewFlatRatePriceCalculator(NewPriceFromUSDcents(8)),
	{City: "Oak City", FractionType: ConstructionWaste, VisitorType: "business"}:  NewFlatRatePriceCalculator(NewPriceFromUSDcents(21)),
}
