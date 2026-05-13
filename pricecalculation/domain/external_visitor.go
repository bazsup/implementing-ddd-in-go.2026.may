package domain

import "errors"

var ErrExternalVisitorTypeUnknown = errors.New("unknown visitor type")

const ExternalVisitorTypePrivate = "private"
const externalVisitorTypeBusiness = "business"

var knownExternalVisitorTypes = map[string]bool{
	ExternalVisitorTypePrivate:  true,
	externalVisitorTypeBusiness: true,
}

type ExternalVisitor struct {
	visitorType string
	id          string
	address     string
	city        string
}

func NewExternalVisitor(visitorType string, id string, address string, city string) (ExternalVisitor, error) {
	if _, ok := knownExternalVisitorTypes[visitorType]; !ok {
		return ExternalVisitor{}, ErrExternalVisitorTypeUnknown
	}
	return ExternalVisitor{visitorType, id, address, city}, nil
}

func (ev ExternalVisitor) Type() string {
	return ev.visitorType
}

func (ev ExternalVisitor) ID() string {
	return ev.id
}

func (ev ExternalVisitor) Address() string {
	return ev.address
}

func (ev ExternalVisitor) City() string {
	return ev.city
}

func (ev ExternalVisitor) Equals(other ExternalVisitor) bool {
	return ev.visitorType == other.visitorType &&
		ev.id == other.id &&
		ev.address == other.address &&
		ev.city == other.city
}
