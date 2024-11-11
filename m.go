package main

import (
	"go-http-server/app"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to run application.", "error", err)
		return
	}
}

func run() error {
	db := make(map[string]string)
	handler := app.NewHandler(db)

	s := http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  1 * time.Minute,
		Addr:         ":3000",
		Handler:      handler,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
