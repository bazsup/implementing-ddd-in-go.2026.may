package domain

import "errors"

var ErrUnknownCustomerType = errors.New("unknown customer type")

const customerTypePrivate = "private"
const customerTypeBusiness = "business"

type Customer interface {
	ID() string
	Type() string
	City() string
}

type PrivateCustomer struct {
	id      string
	address Address
}

func NewPrivateCustomer(id string, address Address) PrivateCustomer {
	return PrivateCustomer{id: id, address: address}
}

func (c PrivateCustomer) ID() string {
	return c.id
}

func (c PrivateCustomer) Type() string {
	return customerTypePrivate
}

func (c PrivateCustomer) City() string {
	return c.address.City()
}

type BusinessCustomer struct {
	id      string
	address Address
	email   string
}

func NewBusinessCustomer(id string, address Address, email string) BusinessCustomer {
	return BusinessCustomer{id: id, address: address, email: email}
}

func (c BusinessCustomer) Email() string {
	return c.email
}

func (c BusinessCustomer) ID() string {
	return c.address.City() + "|" + c.address.Street()
}

func (c BusinessCustomer) Type() string {
	return customerTypeBusiness
}

func (c BusinessCustomer) City() string {
	return c.address.City()
}

func NewCustomer(visitorType, id string, address Address, email string) (Customer, error) {
	switch visitorType {
	case customerTypeBusiness:
		return NewBusinessCustomer(id, address, email), nil
	case customerTypePrivate:
		return NewPrivateCustomer(id, address), nil
	default:
		return nil, ErrUnknownCustomerType
	}
}
