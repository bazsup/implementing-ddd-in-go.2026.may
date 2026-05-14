package domain

import "errors"

var ErrInvalidAddress = errors.New("address must have a non-empty street and city")

type Address struct {
	street string
	city   string
}

func NewAddress(street, city string) (Address, error) {
	if street == "" || city == "" {
		return Address{}, ErrInvalidAddress
	}
	return Address{street: street, city: city}, nil
}

func (a Address) Street() string {
	return a.street
}

func (a Address) City() string {
	return a.city
}

func (a Address) Equals(other Address) bool {
	return a.street == other.street && a.city == other.city
}
