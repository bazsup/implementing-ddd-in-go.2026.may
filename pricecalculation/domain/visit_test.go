package domain_test

import (
	. "implementing-ddd-in-go/pricecalculation/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVisit_WithInvalidDate(t *testing.T) {
	visitor, _ := NewExternalVisitor(ExternalVisitorTypePrivate, "person-1", "Pine Street 1", "Pineville")

	_, err := NewVisit("2025-12-32", visitor)

	assert.EqualError(t, err, ErrInvalidDateForVisit.Error())
}
