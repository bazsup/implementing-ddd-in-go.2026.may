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
	last := vh.visits[len(vh.visits)-1]
	count := 0
	for _, v := range vh.visits {
		if v.Date().Year() == last.Date().Year() && v.Date().Month() == last.Date().Month() {
			count++
		}
	}
	return count
}

func (vh *VisitHistory) Reset() {
	vh.visits = []Visit{}
}
