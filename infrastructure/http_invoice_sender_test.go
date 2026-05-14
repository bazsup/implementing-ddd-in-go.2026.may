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

	err := sut.Send("beavers@dam-building.com", 125.37, "USD")

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

	err := sut.Send("x@y.com", 10.0, "USD")

	assert.EqualError(t, err, ErrFailedToSendInvoice.Error())
}

func TestWhenInvoicingPolicy_SendsForBusinessCustomer(t *testing.T) {
	var sentEmail string
	receiver := WhenInvoicingPolicy(domain.BusinessCustomersRequireInvoice, func(email string, amount float64, currency string) error {
		sentEmail = email
		return nil
	})

	err := receiver(domain.PriceCalculated{CustomerType: "business", Email: "biz@co.com", PriceAmount: 50.0, PriceCurrency: "USD"})

	assert.NoError(t, err)
	assert.Equal(t, "biz@co.com", sentEmail)
}

func TestWhenInvoicingPolicy_SkipsForPrivateCustomer(t *testing.T) {
	callCount := 0
	receiver := WhenInvoicingPolicy(domain.BusinessCustomersRequireInvoice, func(email string, amount float64, currency string) error {
		callCount++
		return nil
	})

	err := receiver(domain.PriceCalculated{CustomerType: "private"})

	assert.NoError(t, err)
	assert.Equal(t, 0, callCount)
}
