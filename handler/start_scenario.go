package handler

import "net/http"

func (h *Handler) StartScenario(w http.ResponseWriter, _ *http.Request) {
	h.logger.Info().Msg("handling scenario request")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("{}")); err != nil {
		h.logger.Error().Err(err).Msg("failed to write response")
	}
}
