package steps

import (
	"context"
	"errors"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-articles-api/api"
	"github.com/ONSdigital/dp-articles-api/config"
	"github.com/ONSdigital/dp-articles-api/service"
	"github.com/ONSdigital/dp-articles-api/service/mock"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type Component struct {
	componenttest.ErrorFeature
	svcList        *service.ExternalServiceList
	svc            *service.Service
	errorChan      chan error
	Config         *config.Config
	HTTPServer     *http.Server
	ServiceRunning bool
	apiFeature     *componenttest.APIFeature
}

func NewComponent() (*Component, error) {

	c := &Component{
		HTTPServer:     &http.Server{},
		errorChan:      make(chan error),
		ServiceRunning: false,
	}

	var err error

	c.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	initMock := &mock.InitialiserMock{
		DoGetHTTPServerFunc:    c.DoGetHTTPServer,
		DoGetHealthCheckFunc:   c.DoGetHealthcheckOk,
		DoGetZebedeeClientFunc: c.DoGetZebedeeClient,
	}

	c.svcList = service.NewServiceList(initMock)

	c.apiFeature = componenttest.NewAPIFeature(c.InitialiseService)

	return c, nil
}

func (c *Component) Reset() *Component {
	c.apiFeature.Reset()
	return c
}

func (c *Component) Close() error {
	if c.svc != nil && c.ServiceRunning {
		c.svc.Close(context.Background())
		c.ServiceRunning = false
	}
	return nil
}

func (c *Component) InitialiseService() (http.Handler, error) {
	var err error
	c.svc, err = service.Run(context.Background(), c.Config, c.svcList, "1", "", "", c.errorChan)
	if err != nil {
		return nil, err
	}

	c.ServiceRunning = true
	return c.HTTPServer.Handler, nil
}

func (c *Component) DoGetHealthcheckOk(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
		HandlerFunc:  func(w http.ResponseWriter, req *http.Request) {},
	}, nil
}

func (c *Component) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	c.HTTPServer.Addr = bindAddr
	c.HTTPServer.Handler = router
	return c.HTTPServer
}

func (c *Component) DoGetZebedeeClient(url string) api.ZebedeeClient {
	return &api.ZebedeeClientMock{
		GetBulletinFunc: func(ctx context.Context, userAccessToken, collectionID, lang, uri string) (zebedee.Bulletin, error) {
			if uri == "/gdp/economy" {
				return zebedee.Bulletin{
					URI:  uri,
					Type: "bulletin",
					Description: zebedee.Description{
						Title:           "Bulletin test",
						Summary:         "Test summary",
						MetaDescription: "Desc",
						Contact: zebedee.Contact{
							Email:     "contact@ons.gov.uk",
							Name:      "Contact name",
							Telephone: "029",
						},
						NationalStatistic: true,
						LatestRelease:     true,
						Keywords:          []string{"economy", "gdp"},
						Edition:           "2020",
						ReleaseDate:       "2020-07-08T23:00:00.000Z",
					},
					Sections: []zebedee.Section{
						{
							Title:    "Section 1",
							Markdown: "Markdown 1",
						},
						{
							Title:    "Section 2",
							Markdown: "Markdown 2",
						},
					},
					Accordion: []zebedee.Section{
						{
							Title:    "Notes",
							Markdown: "Accordion text",
						},
					},
					RelatedBulletins: []zebedee.Link{
						{
							Title: "other bulletin",
							URI:   "bulletin/uri",
						},
					},
					Tables: []zebedee.Figure{
						{
							Title:    "Table 1.1",
							Filename: "table1",
							Version:  "1",
							URI:      "table/1",
						},
					},
				}, nil
			}
			return zebedee.Bulletin{}, errors.New("Unsupported endpoint")
		},
	}
}
