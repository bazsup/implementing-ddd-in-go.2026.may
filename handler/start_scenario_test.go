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
	var reinitCount uint
	stubReinitContext := func ()  { reinitCount++ }

	h := handler.NewHandler(zerolog.Nop(), stubReinitContext, dummyGetVisitorByID)
	h.StartScenario(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{}`, rec.Body.String())
	assert.Equal(t, uint(1), reinitCount)
}
