package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/log.go/v2/log"
)

func LegacyHandler(ctx context.Context, zc ZebedeeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		urlParam := req.URL.Query().Get("url")

		if urlParam == "" {
			err := errors.New("url parameter not found")
			log.Error(ctx, "url parameter not found", err)
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		lang := "en"
		userAccessToken := ""

		response, err := zc.GetBulletin(ctx, userAccessToken, lang, urlParam)
		if err != nil {
			statusCode := http.StatusInternalServerError
			var e zebedee.ErrInvalidZebedeeResponse
			if errors.As(err, &e) {
				statusCode = e.ActualCode
			}
			setStatusCode(ctx, w, statusCode, "retrieving bulletin from Zebedee", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			setStatusCode(ctx, w, http.StatusInternalServerError, "marshalling response failed", err)
			return
		}

		_, err = w.Write(jsonResponse)
		if err != nil {
			setStatusCode(ctx, w, http.StatusInternalServerError, "writing response failed", err)
			return
		}
	}
}

func setStatusCode(ctx context.Context, w http.ResponseWriter, statusCode int, msg string, err error) {
	log.Error(ctx, msg, err)
	w.WriteHeader(statusCode)
}
