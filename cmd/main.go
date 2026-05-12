package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/rs/zerolog"

	"implementing-ddd-in-go/cmd/config"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", statusHandler)

	server := &http.Server{
		Addr:    conf.ServerPort,
		Handler: mux,
	}

	logger.Info().Str("port", conf.ServerPort).Msg("starting server")

	if err = server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func statusHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(map[string]string{"status": "ok"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
