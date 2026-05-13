package domain

import (
	"errors"
	"time"
)

var ErrInvalidDateForVisit = errors.New("invalid date given for Visit")

type Visit struct {
	personId string
	date     time.Time
}

func NewVisit(personId string, date string) (Visit, error) {
	visitDate, dateErr := time.Parse(time.DateOnly, date)
	if dateErr != nil {
		return Visit{}, ErrInvalidDateForVisit
	}
	return Visit{personId: personId, date: visitDate}, nil
}

func (v Visit) inSameMonth(other Visit) bool {
	return v.PersonId() == other.PersonId() &&
		v.Date().Year() == other.Date().Year() &&
		v.Date().Month() == other.Date().Month()
}

func (v Visit) PersonId() string {
	return v.personId
}

func (v Visit) Date() time.Time {
	return v.date
}
