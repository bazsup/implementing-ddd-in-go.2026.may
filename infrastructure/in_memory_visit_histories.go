package infrastructure

import "implementing-ddd-in-go/pricecalculation/domain"

type InMemoryVisitHistories struct {
	histories map[string]*domain.VisitHistory
}

func NewInMemoryVisitHistories() *InMemoryVisitHistories {
	return &InMemoryVisitHistories{histories: make(map[string]*domain.VisitHistory)}
}

func (h *InMemoryVisitHistories) GetByPersonID(personID string) *domain.VisitHistory {
	if vh, ok := h.histories[personID]; ok {
		return vh
	}
	vh := domain.NewVisitHistory(personID)
	h.histories[personID] = vh
	return vh
}

func (h *InMemoryVisitHistories) Save(vh *domain.VisitHistory) {
	h.histories[vh.PersonId()] = vh
}

func (h *InMemoryVisitHistories) Reset() {
	h.histories = make(map[string]*domain.VisitHistory)
}
