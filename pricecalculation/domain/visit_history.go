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
	return 0
}
