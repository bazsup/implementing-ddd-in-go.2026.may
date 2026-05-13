package domain

import "time"

type Visit struct {
	personId string
	date     time.Time
}

func NewVisit(personId string, date time.Time) Visit {
	return Visit{personId: personId, date: date}
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
