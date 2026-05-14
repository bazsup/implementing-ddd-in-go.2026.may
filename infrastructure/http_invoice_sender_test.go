package infrastructure_test

import (
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

func TestSendInvoiceOnPriceCalculated_SendsInvoiceForBusinessCustomer(t *testing.T) {
	var capturedBody []byte
	sut := NewHTTPInvoiceSender("http://stub", func(req *http.Request) (*http.Response, error) {
		capturedBody, _ = io.ReadAll(req.Body)
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}, nil
	}, zerolog.Nop(), domain.BusinessCustomersRequireInvoice)

	event := domain.PriceCalculated{
		CustomerType:  "business",
		PriceAmount:   125.37,
		PriceCurrency: "USD",
		Email:         "beavers@dam-building.com",
	}
	err := sut.SendInvoiceOnPriceCalculated(event)

	assert.NoError(t, err)
	var req struct {
		Email           string  `json:"email"`
		InvoiceAmount   float64 `json:"invoice_amount"`
		InvoiceCurrency string  `json:"invoice_currency"`
	}
	_ = json.Unmarshal(capturedBody, &req)
	assert.Equal(t, "beavers@dam-building.com", req.Email)
	assert.Equal(t, 125.37, req.InvoiceAmount)
	assert.Equal(t, "USD", req.InvoiceCurrency)
}

func TestSendInvoiceOnPriceCalculated_SkipsInvoiceForPrivateCustomer(t *testing.T) {
	requestCount := 0
	sut := NewHTTPInvoiceSender("http://stub", func(req *http.Request) (*http.Response, error) {
		requestCount++
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}, nil
	}, zerolog.Nop(), domain.BusinessCustomersRequireInvoice)

	err := sut.SendInvoiceOnPriceCalculated(domain.PriceCalculated{CustomerType: "private"})

	assert.NoError(t, err)
	assert.Equal(t, 0, requestCount)
}

func TestSendInvoiceOnPriceCalculated_ReturnsErrorOnNonOKStatus(t *testing.T) {
	sut := NewHTTPInvoiceSender("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))}, nil
	}, zerolog.Nop(), domain.BusinessCustomersRequireInvoice)

	err := sut.SendInvoiceOnPriceCalculated(domain.PriceCalculated{CustomerType: "business", Email: "x@y.com"})

	assert.EqualError(t, err, ErrFailedToSendInvoice.Error())
}
