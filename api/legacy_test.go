package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLegacyHandler(t *testing.T) {

	Convey("Given a Legacy handler ", t, func() {
		handler := LegacyHandler(context.Background())

		Convey("when a valid request is received", func() {
			req := httptest.NewRequest("GET", "http://localhost:8080/legacy?url=/gdp/economy", nil)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			Convey("Then the response is correct", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
				So(resp.Body.String(), ShouldResemble, `{"url":"/gdp/economy"}`)
			})
		})

		Convey("when a request without an url parameter is received", func() {
			req := httptest.NewRequest("GET", "http://localhost:8080/legacy", nil)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			Convey("Then an error is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusNotFound)
				So(resp.Body.String(), ShouldResemble, "URL not found\n")
			})
		})

		Convey("when a request with an empty url parameter is received", func() {
			req := httptest.NewRequest("GET", "http://localhost:8080/legacy?url=", nil)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			Convey("Then an error is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusNotFound)
				So(resp.Body.String(), ShouldResemble, "URL not found\n")
			})
		})
	})
}
