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

func TestSend_SendsInvoice(t *testing.T) {
	var capturedBody []byte
	sut := NewHTTPInvoiceSender("http://stub", func(req *http.Request) (*http.Response, error) {
		capturedBody, _ = io.ReadAll(req.Body)
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(""))}, nil
	}, zerolog.Nop())

	event := domain.PriceCalculated{
		PriceAmount:   125.37,
		PriceCurrency: "USD",
		Email:         "beavers@dam-building.com",
	}
	err := sut.Send(event)

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

func TestSend_ReturnsErrorOnNonOKStatus(t *testing.T) {
	sut := NewHTTPInvoiceSender("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))}, nil
	}, zerolog.Nop())

	err := sut.Send(domain.PriceCalculated{Email: "x@y.com"})

	assert.EqualError(t, err, ErrFailedToSendInvoice.Error())
}
