package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
)

type LegacyResponse struct {
	URL string `json:"url,omitempty"`
}

func LegacyHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		urlParam := req.URL.Query().Get("url")

		if urlParam == "" {
			err := errors.New("url parameter not found")
			log.Error(ctx, "url parameter not found", err)
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		response := LegacyResponse{
			URL: urlParam,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Error(ctx, "marshalling response failed", err)
			http.Error(w, "Failed to marshall json response", http.StatusInternalServerError)
			return
		}

		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Error(ctx, "writing response failed", err)
			http.Error(w, "Failed to write http response", http.StatusInternalServerError)
			return
		}
	}
}
