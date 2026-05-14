package main

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"

	"implementing-ddd-in-go/cmd/config"
	"implementing-ddd-in-go/handler"
	"implementing-ddd-in-go/infrastructure"
	"implementing-ddd-in-go/pricecalculation"
	"implementing-ddd-in-go/pricecalculation/domain"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()

	externalVisitors := infrastructure.NewHTTPExternalVisitors(conf.WorkshopServerURL, http.DefaultClient.Do, logger)
	visitHistories := infrastructure.NewInMemoryVisitHistories()
	context := pricecalculation.NewContext(logger, visitHistories.Reset)

	h := handler.NewHandler(
		logger,
		context.Initialize,
		externalVisitors.GetCustomerByPersonID,
		visitHistories.GetByCustomerID,
		visitHistories.Save,
		domain.DefaultFractionPricingPolicy,
		func(domain.PriceCalculated) error { return nil },
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Status)
	mux.HandleFunc("POST /startScenario", h.StartScenario)
	mux.HandleFunc("POST /calculatePrice", h.CalculatePrice)

	server := &http.Server{
		Addr:    conf.ServerPort,
		Handler: mux,
	}

	logger.Info().Str("port", conf.ServerPort).Msg("starting server")

	if err = server.ListenAndServe(); err != nil {
		panic(err)
	}
}
