package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/log.go/v2/log"
)

func LegacyHandler(ctx context.Context, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		handleLegacy(w, r, lang, collectionID, accessToken, zc)
	})
}

func handleLegacy(w http.ResponseWriter, req *http.Request, lang, collectionID, accessToken string, zc ZebedeeClient) {
	ctx := req.Context()

	urlParam := req.URL.Query().Get("url")

	if urlParam == "" {
		err := errors.New("url parameter not found")
		log.Error(ctx, "url parameter not found", err)
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	response, err := zc.GetBulletin(ctx, accessToken, collectionID, lang, urlParam)
	if err != nil {
		setStatusCode(ctx, w, "retrieving bulletin from Zebedee", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		setStatusCode(ctx, w, "marshalling response failed", err)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		setStatusCode(ctx, w, "writing response failed", err)
		return
	}
}

func setStatusCode(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	statusCode := http.StatusInternalServerError
	var e zebedee.ErrInvalidZebedeeResponse
	if errors.As(err, &e) {
		statusCode = e.ActualCode
	}
	log.Error(ctx, msg, err)
	w.WriteHeader(statusCode)
}
