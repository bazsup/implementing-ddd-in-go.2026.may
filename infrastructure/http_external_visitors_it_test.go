//go:build integration
package infrastructure_test

import (
	. "implementing-ddd-in-go/infrastructure"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetVisitorByID_Live(t *testing.T) {
	// arrange
	viper.SetConfigFile("../.env")
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("could not read config: %v", err)
	}
	url := viper.GetString("WORKSHOP_SERVER_URL")
	if url == "" {
		t.Fatal("WORKSHOP_SERVER_URL not set in .env")
	}
	sut := NewHTTPExternalVisitors(url, http.DefaultClient.Do, zerolog.Nop())

	// act
	visitor, err := sut.GetCustomerByPersonID("Squirrel Gus")

	// assert
	assert.NoError(t, err)
	assert.Equal(t, "Squirrel Gus", visitor.ID())
}
