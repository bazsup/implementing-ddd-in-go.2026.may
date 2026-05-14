package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"

	"implementing-ddd-in-go/pricecalculation/domain"
)

var ErrFailedToSendInvoice = errors.New("failed to send invoice")

type invoiceRequest struct {
	Email           string  `json:"email"`
	InvoiceAmount   float64 `json:"invoice_amount"`
	InvoiceCurrency string  `json:"invoice_currency"`
}

type httpInvoiceSender struct {
	baseURL   string
	doRequest forDoingHttpRequest
	logger    zerolog.Logger
}

func NewHTTPInvoiceSender(baseURL string, doRequest forDoingHttpRequest, logger zerolog.Logger) httpInvoiceSender {
	return httpInvoiceSender{baseURL: baseURL, doRequest: doRequest, logger: logger}
}

func WhenInvoicingPolicy(policy domain.InvoicingPolicy, send func(domain.DomainEvent) error) func(domain.DomainEvent) error {
	return func(event domain.DomainEvent) error {
		pc, ok := event.(domain.PriceCalculated)
		if !ok || !policy(pc) {
			return nil
		}
		return send(event)
	}
}

func (s httpInvoiceSender) Send(event domain.DomainEvent) error {
	pc, ok := event.(domain.PriceCalculated)
	if !ok {
		return nil
	}

	payload := invoiceRequest{
		Email:           pc.Email,
		InvoiceAmount:   pc.PriceAmount,
		InvoiceCurrency: pc.PriceCurrency,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error().Err(err).Msg("marshaling invoice request")
		return ErrFailedToSendInvoice
	}

	req, err := http.NewRequest(http.MethodPost, s.baseURL+"/api/invoice", bytes.NewReader(body))
	if err != nil {
		s.logger.Error().Err(err).Msg("creating invoice request")
		return ErrFailedToSendInvoice
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.doRequest(req)
	if err != nil {
		s.logger.Error().Err(err).Msg("executing invoice request")
		return ErrFailedToSendInvoice
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Error().Err(err).Msg("closing invoice response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error().Int("status", resp.StatusCode).Msg("unexpected status from invoice API")
		return ErrFailedToSendInvoice
	}

	return nil
}
