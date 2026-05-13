package domain

type VisitHistory struct {
	personId string
	visits   []Visit
}

func NewVisitHistory(personId string) *VisitHistory {
	return &VisitHistory{personId: personId}
}

func (vh *VisitHistory) PersonId() string {
	return vh.personId
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

func (vh *VisitHistory) CalculatePriceOfVisit(visit Visit, droppedFractions []DroppedFraction) Price {
	vh.Add(visit)
	var total Price
	for _, df := range droppedFractions {
		total = total.Add(df.CalculatePrice())
	}
	if vh.NumberOfVisitsInMonthOfLastVisit() >= 3 {
		total = total.AddFee(5)
	}
	return total
}

func (vh *VisitHistory) Reset() {
	vh.visits = []Visit{}
}
