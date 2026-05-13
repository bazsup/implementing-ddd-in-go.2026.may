package domain_test

import (
	. "implementing-ddd-in-go/pricecalculation/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVisit_WithInvalidDate(t *testing.T) {
	// act
	_, err := NewVisit("Bald Eagle", "2025-12-32")

	// assert
	assert.EqualError(t, err, ErrInvalidDateForVisit.Error())
}
