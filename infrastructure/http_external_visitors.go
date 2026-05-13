package infrastructure

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog"

	"implementing-ddd-in-go/pricecalculation/domain"
)

var ErrFailedToGetVisitors = errors.New("failed to get external visitors")
var ErrExternalVisitorNotFound = errors.New("external visitor not found")

type forDoingHttpRequest func(*http.Request) (*http.Response, error)

type httpExternalVisitors struct {
	baseURL   string
	doRequest forDoingHttpRequest
	logger    zerolog.Logger
}

func NewHTTPExternalVisitors(baseURL string, doRequest forDoingHttpRequest, logger zerolog.Logger) httpExternalVisitors {
	return httpExternalVisitors{baseURL: baseURL, doRequest: doRequest, logger: logger}
}

type ExternalVisitorResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Address string `json:"address"`
	City    string `json:"city"`
	Email   string `json:"email"`
}

func (h httpExternalVisitors) GetVisitorByID(personID string) (domain.ExternalVisitor, error) {
	req, err := http.NewRequest(http.MethodGet, h.baseURL+"/api/users", nil)
	if err != nil {
		h.logger.Error().Err(err).Msg("creating request")
		return domain.ExternalVisitor{}, ErrFailedToGetVisitors
	}

	resp, err := h.doRequest(req)
	if err != nil {
		h.logger.Error().Err(err).Msg("executing request")
		return domain.ExternalVisitor{}, ErrFailedToGetVisitors
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			h.logger.Error().Err(err).Msg("closing response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		h.logger.Error().Int("status", resp.StatusCode).Msg("unexpected status")
		return domain.ExternalVisitor{}, ErrFailedToGetVisitors
	}

	var visitors []ExternalVisitorResponse
	if err := json.NewDecoder(resp.Body).Decode(&visitors); err != nil {
		h.logger.Error().Err(err).Msg("decoding response")
		return domain.ExternalVisitor{}, ErrFailedToGetVisitors
	}

	for _, v := range visitors {
		if v.ID == personID {
			return domain.NewExternalVisitor(v.Type, v.ID, v.Address, v.City)
		}
	}

	return domain.ExternalVisitor{}, ErrExternalVisitorNotFound
}
