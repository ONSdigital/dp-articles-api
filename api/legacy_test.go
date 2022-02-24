package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/headers"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLegacyHandler(t *testing.T) {

	Convey("Given a Legacy handler ", t, func() {
		url := "/gdp/economy"
		b := zebedee.Bulletin{URI: url, Description: zebedee.Description{Title: "Test"}}
		mockZebedeeClient := &ZebedeeClientMock{
			GetBulletinFunc: func(ctx context.Context, userAccessToken, collectionId, lang, uri string) (zebedee.Bulletin, error) {
				return b, nil
			},
		}
		handler := LegacyHandler(context.Background(), mockZebedeeClient)

		Convey("when a valid request without headers is received (web)", func() {
			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8080/v1/articles/legacy?url=%s", url), nil)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			Convey("Then the call to Zebedee is correct", func() {
				So(len(mockZebedeeClient.GetBulletinCalls()), ShouldEqual, 1)
				So(mockZebedeeClient.GetBulletinCalls()[0].UserAccessToken, ShouldEqual, "")
				So(mockZebedeeClient.GetBulletinCalls()[0].CollectionID, ShouldEqual, "")
				So(mockZebedeeClient.GetBulletinCalls()[0].Lang, ShouldEqual, "en")
				So(mockZebedeeClient.GetBulletinCalls()[0].URI, ShouldEqual, url)
			})
			Convey("And the response is correct", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
				expectedJson, _ := json.Marshal(b)
				So(resp.Body.Bytes(), ShouldResemble, expectedJson)
				So(len(mockZebedeeClient.GetBulletinCalls()), ShouldEqual, 1)
			})
		})

		Convey("when a valid request with headers is received (publishing)", func() {
			accessToken := "user-access-token"
			collectionID := "my-collection"

			req := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:8080/v1/articles/legacy?url=%s", url), nil)
			headers.SetAuthToken(req, accessToken)
			headers.SetCollectionID(req, collectionID)

			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			Convey("Then the call to Zebedee is correct", func() {
				So(len(mockZebedeeClient.GetBulletinCalls()), ShouldEqual, 1)
				So(mockZebedeeClient.GetBulletinCalls()[0].UserAccessToken, ShouldEqual, accessToken)
				So(mockZebedeeClient.GetBulletinCalls()[0].CollectionID, ShouldEqual, collectionID)
				So(mockZebedeeClient.GetBulletinCalls()[0].Lang, ShouldEqual, "en")
				So(mockZebedeeClient.GetBulletinCalls()[0].URI, ShouldEqual, url)
			})
			Convey("And the response is correct", func() {
				So(resp.Code, ShouldEqual, http.StatusOK)
				expectedJson, _ := json.Marshal(b)
				So(resp.Body.Bytes(), ShouldResemble, expectedJson)
				So(len(mockZebedeeClient.GetBulletinCalls()), ShouldEqual, 1)
			})
		})

		Convey("when a request without an url parameter is received", func() {
			req := httptest.NewRequest("GET", "http://localhost:8080/v1/articles/legacy", nil)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			Convey("Then an error is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusNotFound)
				So(resp.Body.String(), ShouldResemble, "URL not found\n")
			})
		})

		Convey("when a request with an empty url parameter is received", func() {
			req := httptest.NewRequest("GET", "http://localhost:8080/v1/articles/legacy?url=", nil)
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			Convey("Then an error is returned", func() {
				So(resp.Code, ShouldEqual, http.StatusNotFound)
				So(resp.Body.String(), ShouldResemble, "URL not found\n")
			})
		})
	})
}
