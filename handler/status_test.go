package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"implementing-ddd-in-go/handler"
)

func TestStatusHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h := handler.NewHandler(zerolog.Nop(), dummyReinitContext, dummyGetVisitorByID, dummyGetVisitHistoryByPersonID, dummySaveVisitHistory)
	h.Status(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}
