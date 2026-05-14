package infrastructure_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	. "implementing-ddd-in-go/infrastructure"
	"implementing-ddd-in-go/pricecalculation/domain"
)

func TestGetCustomerByPersonID_ReturnsPrivateCustomer(t *testing.T) {
	body, _ := json.Marshal([]ExternalVisitorResponse{
		{ID: "Squirrel Gus", Type: "private", Address: "Lowest Branch 19", City: "Oak City"},
	})
	sut := NewHTTPExternalVisitors("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(body))}, nil
	}, zerolog.Nop())

	customer, err := sut.GetCustomerByPersonID("Squirrel Gus")

	assert.NoError(t, err)
	assert.Equal(t, "Squirrel Gus", customer.ID())
	assert.Equal(t, "private", customer.Type())
	assert.Equal(t, "Oak City", customer.City())
}

func TestGetCustomerByPersonID_ReturnsBusinessCustomerWithAddressDerivedID(t *testing.T) {
	body, _ := json.Marshal([]ExternalVisitorResponse{
		{ID: "Bear Billy", Type: "business", Address: "Trunk 35", City: "Pineville"},
	})
	sut := NewHTTPExternalVisitors("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(body))}, nil
	}, zerolog.Nop())

	customer, err := sut.GetCustomerByPersonID("Bear Billy")

	addr, _ := domain.NewAddress("Trunk 35", "Pineville")
	expected := domain.NewBusinessCustomer("Bear Billy", addr, "")
	assert.NoError(t, err)
	assert.Equal(t, expected.ID(), customer.ID())
	assert.Equal(t, "business", customer.Type())
}

func TestGetCustomerByPersonID_TwoEmployeesAtSameAddressHaveSameID(t *testing.T) {
	body, _ := json.Marshal([]ExternalVisitorResponse{
		{ID: "employee-a", Type: "business", Address: "Oak Avenue 1", City: "Oak City"},
		{ID: "employee-b", Type: "business", Address: "Oak Avenue 1", City: "Oak City"},
	})
	sut := NewHTTPExternalVisitors("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(body))}, nil
	}, zerolog.Nop())

	customerA, errA := sut.GetCustomerByPersonID("employee-a")
	customerB, errB := sut.GetCustomerByPersonID("employee-b")

	assert.NoError(t, errA)
	assert.NoError(t, errB)
	assert.Equal(t, customerA.ID(), customerB.ID())
}

func TestGetCustomerByPersonID_ReturnsBusinessCustomerWithEmail(t *testing.T) {
	body, _ := json.Marshal([]ExternalVisitorResponse{
		{ID: "Beaver Bob", Type: "business", Address: "Dam Road 1", City: "Oak City", Email: "bob@dam.com"},
	})
	sut := NewHTTPExternalVisitors("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(body))}, nil
	}, zerolog.Nop())

	customer, err := sut.GetCustomerByPersonID("Beaver Bob")

	assert.NoError(t, err)
	bc, ok := customer.(domain.BusinessCustomer)
	assert.True(t, ok)
	assert.Equal(t, "bob@dam.com", bc.Email())
}

func TestGetCustomerByPersonID_VisitorNotFound(t *testing.T) {
	sut := NewHTTPExternalVisitors("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`[]`))}, nil
	}, zerolog.Nop())

	_, err := sut.GetCustomerByPersonID("Squirrel Gus")

	assert.EqualError(t, err, ErrExternalVisitorNotFound.Error())
}

func TestGetCustomerByPersonID_NonOKStatus(t *testing.T) {
	sut := NewHTTPExternalVisitors("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))}, nil
	}, zerolog.Nop())

	_, err := sut.GetCustomerByPersonID("Squirrel Gus")

	assert.EqualError(t, err, ErrFailedToGetVisitors.Error())
}

func TestGetCustomerByPersonID_InvalidJSON(t *testing.T) {
	sut := NewHTTPExternalVisitors("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`not json`))}, nil
	}, zerolog.Nop())

	_, err := sut.GetCustomerByPersonID("Squirrel Gus")

	assert.EqualError(t, err, ErrFailedToGetVisitors.Error())
}
