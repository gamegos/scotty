# Scotty—push notification service

**WORK IN PROGRESS. Things are likely to change.**

## Installation

Install binary with go get:
```
go get github.com/gamegos/scotty
```

## Running
To start the HTTP server and accept requests to the API, run:
```
scotty -config=/path/to/config.toml
```
By default, `config.toml` file in the same directory as the binary will be used as configuration.

## Tests
```
go test -v ./...
```

## License
MIT. See [LICENSE](./LICENSE).
