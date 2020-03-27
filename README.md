# Sum Metric Service Problem

Metric logging and reporting service that sums metrics by time window for the most recent hour.

## Requirements

- Go 1.14+

## Getting Started

1. Download dependencies `go mod tidy`
2. Start webserver: `go run .`. It starts on port `9000`
3. Call the apis via curl or a http client:
   - `GET localhost:9000/metric/{key}/sum`
   - `POST localhost:9000/metric/{key}`

To run all tests use `go test -v ./...` command.

To build a binary, use `go build ./...` command.

## Notes

1. Sum for each metric key is stored in-memory map. Each log event
is not stored because the sum calculation can happen without it
2. Using the standard library api to remove events that are
 older than one hour: `time.AfterFunc`
3. Webserver port and TTL time to remove older events are hard coded
4. Everything is done via the standard library except the
routes matching. Using `github.com/julienschmidt/httprouter` package for it
5. Everything is one, `main`, package because there are only 3 files.
Hence, no method/function/struct were exported
6. There is some code to gracefully shutdown the webserver
