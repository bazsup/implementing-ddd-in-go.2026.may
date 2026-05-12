package handler

import (
	"encoding/json"
	"net/http"
)

type calculatePriceRequest struct {
	PersonID         string            `json:"person_id"`
	VisitID          string            `json:"visit_id"`
	Date             string            `json:"date"`
	DroppedFractions []droppedFraction `json:"dropped_fractions"`
}

type droppedFraction struct {
	AmountDropped int    `json:"amount_dropped"`
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := calculatePriceResponse{
		PersonID:      req.PersonID,
		VisitID:       req.VisitID,
		PriceAmount:   1,
		PriceCurrency: "USD",
	}

	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
