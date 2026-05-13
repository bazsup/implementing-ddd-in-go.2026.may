package handler

import (
	"encoding/json"
	"net/http"

	"implementing-ddd-in-go/pricecalculation"
)

type calculatePriceRequest struct {
	PersonID         string            `json:"person_id"`
	VisitID          string            `json:"visit_id"`
	Date             string            `json:"date"`
	DroppedFractions []droppedFraction `json:"dropped_fractions"`
}

type droppedFraction struct {
	AmountDropped uint   `json:"amount_dropped"`
	FractionType  string `json:"fraction_type"`
}

type calculatePriceResponse struct {
	PersonID      string  `json:"person_id"`
	VisitID       string  `json:"visit_id"`
	PriceAmount   float64 `json:"price_amount"`
	PriceCurrency string  `json:"price_currency"`
}

func (h *Handler) CalculatePrice(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("handling calculatePrice request")

	var req calculatePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Msg("failed to decode request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var fractions []pricecalculation.RawDroppedFraction
	for _, f := range req.DroppedFractions {
		fractions = append(fractions, pricecalculation.RawDroppedFraction{Type: f.FractionType, AmountKG: f.AmountDropped})
	}

	priceCalculator := pricecalculation.NewPriceCalculator(h.getVisitorByID, h.getVisitHistoryByPersonID, h.saveVisitHistory)
	calculatedPrice, err := priceCalculator.CalculatePrice(req.PersonID, req.VisitID, req.Date, fractions)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to calculate price")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := calculatePriceResponse{
		PersonID:      calculatedPrice.PersonID,
		VisitID:       calculatedPrice.VisitID,
		PriceAmount:   calculatedPrice.PriceAmount,
		PriceCurrency: calculatedPrice.PriceCurrency,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to marshal response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(body); err != nil {
		h.logger.Error().Err(err).Msg("failed to write response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
