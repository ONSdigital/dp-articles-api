
---

:warning: This repository will be archived in February 2023 as it is not currently in development. :warning:

---

# dp-articles-api
Articles API for ONS Digital Publishing

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default                   | Description
| ---------------------------- | ------------------------- | -----------
| BIND_ADDR                    | :27000                    | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                        | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s                       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| API_ROUTER_URL               | http://localhost:23200/v1 | The URL of the [dp-api-router](https://github.com/ONSdigital/dp-api-router)

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

