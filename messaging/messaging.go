package messaging

import "implementing-ddd-in-go/pricecalculation/domain"

type MessageSender func(domain.DomainEvent) error

type MessageReceiver func(domain.DomainEvent) error

type Bus struct {
	receivers map[string][]MessageReceiver
}

func NewBus() *Bus {
	return &Bus{receivers: make(map[string][]MessageReceiver)}
}

func (b *Bus) Register(eventType string, receiver MessageReceiver) {
	b.receivers[eventType] = append(b.receivers[eventType], receiver)
}

func (b *Bus) Send(event domain.DomainEvent) error {
	for _, r := range b.receivers[event.Type()] {
		if err := r(event); err != nil {
			return err
		}
	}
	return nil
}
