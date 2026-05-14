package domain

type DomainEvent interface {
	Type() string
}
