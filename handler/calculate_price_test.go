package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"implementing-ddd-in-go/handler"
)

func TestCalculatePriceHandler(t *testing.T) {
	body := `{
		"date": "2023-07-23",
		"dropped_fractions": [
			{"amount_dropped": 15, "fraction_type": "Green waste"},
			{"amount_dropped": 39, "fraction_type": "Construction waste"}
		],
		"person_id": "Bald Eagle",
		"visit_id": "1"
	}`

	req := httptest.NewRequest(http.MethodPost, "/calculatePrice", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h := handler.NewHandler(zerolog.Nop(), dummyReinitContext, dummyGetVisitorByID, dummyGetVisitHistoryByPersonID, dummySaveVisitHistory, dummyFractionPricingPolicy, dummyPublishDomainEvent)
	h.CalculatePrice(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		PersonID      string  `json:"person_id"`
		VisitID       string  `json:"visit_id"`
		PriceAmount   float64 `json:"price_amount"`
		PriceCurrency string  `json:"price_currency"`
	}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "Bald Eagle", resp.PersonID)
	assert.Equal(t, "1", resp.VisitID)
	assert.NotZero(t, resp.PriceAmount)
	assert.Equal(t, "USD", resp.PriceCurrency)
}
