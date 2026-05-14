package domain

import (
	"errors"
	"time"
)

var ErrInvalidDateForVisit = errors.New("invalid date given for Visit")

type Visit struct {
	date      time.Time
	visitor   Customer
	fractions []DroppedFraction
}

func NewVisit(date string, visitor Customer) (Visit, error) {
	visitDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return Visit{}, ErrInvalidDateForVisit
	}
	return Visit{date: visitDate, visitor: visitor}, nil
}

func (v Visit) inSameMonth(other Visit) bool {
	return v.visitor.ID() == other.visitor.ID() &&
		v.Date().Year() == other.Date().Year() &&
		v.Date().Month() == other.Date().Month()
}

func (v Visit) Date() time.Time {
	return v.date
}

func (v Visit) Visitor() Customer {
	return v.visitor
}

func (v Visit) withDroppedFractions(fractions []DroppedFraction) Visit {
	return Visit{date: v.date, visitor: v.visitor, fractions: fractions}
}
