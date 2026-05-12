package handler

import "net/http"

func (h *Handler) StartScenario(w http.ResponseWriter, _ *http.Request) {
	h.logger.Info().Msg("handling scenario request")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
