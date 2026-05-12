package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) Status(w http.ResponseWriter, _ *http.Request) {
	h.logger.Info().Msg("handling status request")

	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(map[string]string{"status": "ok"})
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to marshal response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(body); err != nil {
		h.logger.Error().Err(err).Msg("failed to write response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
