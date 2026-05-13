package domain

type VisitHistory struct {
	visits []Visit
}

func NewVisitHistory() *VisitHistory {
	return &VisitHistory{}
}

func (vh *VisitHistory) Add(visit Visit) {
	vh.visits = append(vh.visits, visit)
}

func (vh *VisitHistory) NumberOfVisitsInMonthOfLastVisit() int {
	if len(vh.visits) == 0 {
		return 0
	}
	count := 0
	for _, v := range vh.visits {
		if v.inSameMonth(vh.lastVisit()) {
			count++
		}
	}
	return count
}

func (vh *VisitHistory) lastVisit() Visit {
	if len(vh.visits) == 0 {
		return Visit{}
	}

	return vh.visits[len(vh.visits)-1]
}

func (vh *VisitHistory) Reset() {
	vh.visits = []Visit{}
}
