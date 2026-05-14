package domain_test

import (
	. "implementing-ddd-in-go/pricecalculation/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVisit_WithInvalidDate(t *testing.T) {
	addr, _ := NewAddress("Pine Street 1", "Pineville")
	customer, _ := NewCustomer("private", "person-1", addr, "")

	_, err := NewVisit("2025-12-32", customer)

	assert.EqualError(t, err, ErrInvalidDateForVisit.Error())
}
