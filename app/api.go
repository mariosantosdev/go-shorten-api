package app

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(db map[string]string) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Post("/shorten", handleShortUrl(db))
	r.Get("/{code}", handleGetUrl(db))

	return r
}

func sendJSON(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)

	if err != nil {
		slog.Error("fail to marshal json", "error", err)
		sendJSON(
			w,
			Response{Error: "something went wrong"},
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("fail to send json data", "error", err)
		return
	}
}

type PostBody struct {
	URL string `json:"url"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func handleShortUrl(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body PostBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJSON(
				w,
				Response{Error: "invalid body"},
				http.StatusUnprocessableEntity,
			)
			return
		}

		if _, err := url.Parse(body.URL); err != nil {
			sendJSON(
				w,
				Response{Error: "invalid url"},
				http.StatusBadRequest,
			)
			return
		}

		code := strconv.FormatInt(rand.Int63(), 10)
		db[code] = body.URL

		sendJSON(
			w,
			Response{Data: code},
			http.StatusCreated,
		)
	}
}

func handleGetUrl(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		data, ok := db[code]

		if !ok {
			http.Error(w, "url does not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, data, http.StatusPermanentRedirect)
	}
}
