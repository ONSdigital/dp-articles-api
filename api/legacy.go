package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
)

// ClientError is an interface that can be used to retrieve the status code if a client has errored
type ClientError interface {
	Error() string
	Code() int
}

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
			setStatusCode(req, w, "retrieving bulletin from Zebedee", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			setStatusCode(req, w, "marshalling response failed", err)
			return
		}

		_, err = w.Write(jsonResponse)
		if err != nil {
			setStatusCode(req, w, "writing response failed", err)
			return
		}
	}
}

func setStatusCode(req *http.Request, w http.ResponseWriter, msg string, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}
	log.Error(req.Context(), msg, err)
	w.WriteHeader(status)
}
