package domain

import "errors"

var ErrConcurrentModification = errors.New("visit history was modified by another request")

type VisitHistory struct {
	customerID string
	visits     []Visit
	version    int
}

func NewVisitHistory(customerID string) *VisitHistory {
	return &VisitHistory{customerID: customerID}
}

func HydrateVisitHistory(customerID string, visits []Visit, version int) *VisitHistory {
	visitsCopy := make([]Visit, len(visits))
	copy(visitsCopy, visits)
	return &VisitHistory{customerID: customerID, visits: visitsCopy, version: version}
}

func (vh *VisitHistory) Version() int {
	return vh.version
}

func (vh *VisitHistory) Visits() []Visit {
	return vh.visits
}

func (vh *VisitHistory) CustomerID() string {
	return vh.customerID
}

func (vh *VisitHistory) Add(visit Visit) {
	vh.visits = append(vh.visits, visit)
}

func (vh *VisitHistory) NumberOfVisitsInMonthOfLastVisit() int {
	if len(vh.visits) == 0 {
		return 0
	}
	numberOfVisits := 0
	for _, v := range vh.visits {
		if v.inSameMonth(vh.lastVisit()) {
			numberOfVisits++
		}
	}
	return numberOfVisits
}

func (vh *VisitHistory) lastVisit() Visit {
	if len(vh.visits) == 0 {
		return Visit{}
	}

	return vh.visits[len(vh.visits)-1]
}

func (vh *VisitHistory) AccumulatedWeightForFractionTypeInYear(fractionTypeName string, year int) Weight {
	var total uint
	for _, v := range vh.visits {
		if v.date.Year() != year {
			continue
		}
		for _, f := range v.fractions {
			if f.fractionType.name == fractionTypeName {
				total += f.weight.Amount()
			}
		}
	}
	return NewWeightFromKG(total)
}

func (vh *VisitHistory) CalculatePriceOfVisit(visit Visit, droppedFractions []DroppedFraction, feePolicy FeePolicy, fractionPricingPolicy FractionPricingPolicy) (Price, error) {
	var total Price
	for _, df := range droppedFractions {
		alreadyDropped := vh.AccumulatedWeightForFractionTypeInYear(df.fractionType.name, visit.date.Year())
		calc, err := fractionPricingPolicy.CalculatorFor(df.fractionType.name, visit.visitor.City(), visit.visitor.Type(), alreadyDropped)
		if err != nil {
			return Price{}, err
		}
		total = total.Add(calc.CalculatePrice(df))
	}
	vh.Add(visit.withDroppedFractions(droppedFractions))
	vh.version++
	return feePolicy.AddFee(vh, total), nil
}

func (vh *VisitHistory) Reset() {
	vh.visits = []Visit{}
}
