package domain_test

import (
	. "implementing-ddd-in-go/pricecalculation/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExternalVisitor_TypePrivate(t *testing.T) {
	// arrange
	visitorType := "private"
	id := "Squirrel Gus"
	address := "Lowest Branch 19"
	city := "Oak City"

	// act
	visitor, err := NewExternalVisitor(visitorType, id, address, city)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, visitorType, visitor.Type())
	assert.Equal(t, id, visitor.ID())
	assert.Equal(t, address, visitor.Address())
	assert.Equal(t, city, visitor.City())
}

func TestNewExternalVisitor_TypeBusiness(t *testing.T) {
	// arrange
	visitorType := "business"
	id := "Bear Billy"
	address := "Trunk 35"
	city := "Pineville"

	// act
	visitor, err := NewExternalVisitor(visitorType, id, address, city)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, visitorType, visitor.Type())
	assert.Equal(t, id, visitor.ID())
	assert.Equal(t, address, visitor.Address())
	assert.Equal(t, city, visitor.City())
}

func TestNewExternalVisitor_UnknownType(t *testing.T) {
	// act
	_, err := NewExternalVisitor("unknown", "Bear Billy", "Trunk 35", "Pineville")

	// assert
	assert.EqualError(t, err, ErrExternalVisitorTypeUnknown.Error())
}

func TestExternalVisitors_Equals(t *testing.T) {
	// arrange
	visitor1, _ := NewExternalVisitor("private", "Squirrel Gus", "Lowest Branch 19", "Oak City")
	visitor2, _ := NewExternalVisitor("private", "Squirrel Gus", "Lowest Branch 19", "Oak City")

	// act/assert
	assert.True(t, visitor1.Equals(visitor2))
}
