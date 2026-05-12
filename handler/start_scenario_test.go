package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"implementing-ddd-in-go/handler"
)

func TestStartScenarioHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/startScenario", nil)
	rec := httptest.NewRecorder()

	h := handler.New(zerolog.Nop())
	h.StartScenario(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{}`, rec.Body.String())
}
