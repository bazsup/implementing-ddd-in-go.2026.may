package infrastructure

import "implementing-ddd-in-go/pricecalculation/domain"

type InMemoryVisitHistories struct {
	histories map[string]*domain.VisitHistory
}

func NewInMemoryVisitHistories() *InMemoryVisitHistories {
	return &InMemoryVisitHistories{histories: make(map[string]*domain.VisitHistory)}
}

func (h *InMemoryVisitHistories) GetByPersonID(personID string) *domain.VisitHistory {
	existing, ok := h.histories[personID]
	if !ok {
		existing = domain.NewVisitHistory(personID)
		h.histories[personID] = existing
	}
	return domain.HydrateVisitHistory(personID, existing.Visits(), existing.Version())
}

func (h *InMemoryVisitHistories) Save(vh *domain.VisitHistory) error {
	existing, ok := h.histories[vh.PersonId()]
	if ok && existing.Version()+1 != vh.Version() {
		return domain.ErrConcurrentModification
	}
	h.histories[vh.PersonId()] = domain.HydrateVisitHistory(vh.PersonId(), vh.Visits(), vh.Version())
	return nil
}

func (h *InMemoryVisitHistories) Reset() {
	h.histories = make(map[string]*domain.VisitHistory)
}
