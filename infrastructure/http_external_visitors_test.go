package infrastructure_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	. "implementing-ddd-in-go/infrastructure"
	"implementing-ddd-in-go/pricecalculation/domain"
)

func TestGetVisitorByID_ReturnsVisitor(t *testing.T) {
	// arrange
	body, _ := json.Marshal([]ExternalVisitorResponse{
		{ID: "Squirrel Gus", Type: "private", Address: "Lowest Branch 19", City: "Oak City"},
	})
	sut := NewHTTPExternalVisitors("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(body))}, nil
	}, zerolog.Nop())

	// act
	visitor, err := sut.GetVisitorByID("Squirrel Gus")

	// assert
	expected, _ := domain.NewExternalVisitor("private", "Squirrel Gus", "Lowest Branch 19", "Oak City")
	assert.NoError(t, err)
	assert.True(t, visitor.Equals(expected))
}

func TestGetVisitorByID_VisitorNotFound(t *testing.T) {
	// arrange
	sut := NewHTTPExternalVisitors("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`[]`))}, nil
	}, zerolog.Nop())

	// act
	_, err := sut.GetVisitorByID("Squirrel Gus")

	// assert
	assert.EqualError(t, err, ErrExternalVisitorNotFound.Error())
}

func TestGetVisitorByID_NonOKStatus(t *testing.T) {
	// arrange
	sut := NewHTTPExternalVisitors("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(strings.NewReader(""))}, nil
	}, zerolog.Nop())

	// act
	_, err := sut.GetVisitorByID("Squirrel Gus")

	// assert
	assert.EqualError(t, err, ErrFailedToGetVisitors.Error())
}

func TestGetVisitorByID_InvalidJSON(t *testing.T) {
	// arrange
	sut := NewHTTPExternalVisitors("http://stub", func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`not json`))}, nil
	}, zerolog.Nop())

	// act
	_, err := sut.GetVisitorByID("Squirrel Gus")

	// assert
	assert.EqualError(t, err, ErrFailedToGetVisitors.Error())
}
