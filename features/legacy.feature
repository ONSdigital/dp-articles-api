Feature: Legacy endpoint

  Scenario: Call the legacy endpoint without an url parameter
    When I GET "/articles/legacy"
    Then the HTTP status code should be "404"
    And I should receive the following response:
    """
    URL not found
    """

 Scenario: Call the legacy endpoint with a valid url parameter
    When I GET "/articles/legacy?url=/gdp/economy"
    And I should receive the following JSON response with status "200":
        """
        {
            "accordion": [
              {
              "title":"Notes",
              "markdown":"Accordion text"
              }
            ], 
            "alerts":null, 
            "charts":null, 
            "description": {
              "contact": {
                "email":"contact@ons.gov.uk", 
                "name":"Contact name", 
                "telephone":"029"
              }, 
              "datasetId":"", 
              "edition":"2020", 
              "keywords": ["economy", "gdp"], 
              "latestRelease":true, 
              "metaDescription":"Desc", 
              "nationalStatistic":true, 
              "nextRelease":"", 
              "preUnit":"", 
              "releaseDate":"2020-07-08T23:00:00.000Z", 
              "source":"", 
              "summary":"Test summary", 
              "title":"Bulletin test", 
              "unit":"", 
              "versionLabel":""
            }, 
            "equations":null, 
            "images":null, 
            "latestReleaseUri":"", 
            "links":null, 
            "relatedBulletins":[
              {
              "title":"other bulletin",
              "uri":"bulletin/uri"
              }
            ],
            "relatedData":null, 
            "sections": [
              {
                "title":"Section 1",
                "markdown":"Markdown 1"
              },
              {
                "title":"Section 2",
                "markdown":"Markdown 2"
              }
            ], 
            "tables":[
              {
              "title":"Table 1.1",
              "filename":"table1",
              "version":"1",
              "uri":"table/1"
              }
            ], 
            "type":"bulletin", 
            "uri":"/gdp/economy", 
            "versions":null
        }
        """

  Scenario: Call the legacy endpoint with an invalid url parameter
    When I GET "/articles/legacy?url=/gdp/invalid"
    Then the HTTP status code should be "500"