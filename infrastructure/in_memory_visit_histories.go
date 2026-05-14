package infrastructure

import "implementing-ddd-in-go/pricecalculation/domain"

type InMemoryVisitHistories struct {
	histories map[string]*domain.VisitHistory
}

func NewInMemoryVisitHistories() *InMemoryVisitHistories {
	return &InMemoryVisitHistories{histories: make(map[string]*domain.VisitHistory)}
}

func (h *InMemoryVisitHistories) GetByCustomerID(customerID string) *domain.VisitHistory {
	existing, ok := h.histories[customerID]
	if !ok {
		existing = domain.NewVisitHistory(customerID)
		h.histories[customerID] = existing
	}
	return domain.HydrateVisitHistory(customerID, existing.Visits(), existing.Version())
}

func (h *InMemoryVisitHistories) Save(vh *domain.VisitHistory) error {
	existing, ok := h.histories[vh.CustomerID()]
	if ok && existing.Version()+1 != vh.Version() {
		return domain.ErrConcurrentModification
	}
	h.histories[vh.CustomerID()] = domain.HydrateVisitHistory(vh.CustomerID(), vh.Visits(), vh.Version())
	return nil
}

func (h *InMemoryVisitHistories) Reset() {
	h.histories = make(map[string]*domain.VisitHistory)
}
