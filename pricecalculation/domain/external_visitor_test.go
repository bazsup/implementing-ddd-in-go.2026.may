package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "implementing-ddd-in-go/pricecalculation/domain"
)

func TestNewAddress_ValidStreetAndCity(t *testing.T) {
	addr, err := NewAddress("Lowest Branch 19", "Oak City")

	assert.NoError(t, err)
	assert.Equal(t, "Lowest Branch 19", addr.Street())
	assert.Equal(t, "Oak City", addr.City())
}

func TestNewAddress_EmptyStreet(t *testing.T) {
	_, err := NewAddress("", "Oak City")

	assert.EqualError(t, err, ErrInvalidAddress.Error())
}

func TestNewAddress_EmptyCity(t *testing.T) {
	_, err := NewAddress("Lowest Branch 19", "")

	assert.EqualError(t, err, ErrInvalidAddress.Error())
}

func TestNewCustomer_Private(t *testing.T) {
	addr, _ := NewAddress("Lowest Branch 19", "Oak City")

	customer, err := NewCustomer("private", "Squirrel Gus", addr)

	assert.NoError(t, err)
	assert.Equal(t, "Squirrel Gus", customer.ID())
	assert.Equal(t, "private", customer.Type())
	assert.Equal(t, "Oak City", customer.City())
}

func TestNewCustomer_Business(t *testing.T) {
	addr, _ := NewAddress("Trunk 35", "Pineville")

	customer, err := NewCustomer("business", "Bear Billy", addr)

	assert.NoError(t, err)
	assert.Equal(t, "Pineville|Trunk 35", customer.ID())
	assert.Equal(t, "business", customer.Type())
	assert.Equal(t, "Pineville", customer.City())
}

func TestNewCustomer_UnknownType(t *testing.T) {
	addr, _ := NewAddress("Trunk 35", "Pineville")

	_, err := NewCustomer("unknown", "Bear Billy", addr)

	assert.EqualError(t, err, ErrUnknownCustomerType.Error())
}

func TestBusinessCustomer_TwoEmployeesAtSameAddressHaveSameID(t *testing.T) {
	addr, _ := NewAddress("Oak Avenue 1", "Oak City")

	employeeA, _ := NewCustomer("business", "employee-a", addr)
	employeeB, _ := NewCustomer("business", "employee-b", addr)

	assert.Equal(t, employeeA.ID(), employeeB.ID())
}

func TestBusinessCustomer_DifferentAddressHasDifferentID(t *testing.T) {
	addr1, _ := NewAddress("Oak Avenue 1", "Oak City")
	addr2, _ := NewAddress("Pine Street 1", "Oak City")

	business1, _ := NewCustomer("business", "emp-a", addr1)
	business2, _ := NewCustomer("business", "emp-b", addr2)

	assert.NotEqual(t, business1.ID(), business2.ID())
}
