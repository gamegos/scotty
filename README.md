# Push Service

## Installation

Install binary with go get:
```
go get gitlab.fixb.com/mir/push
```

## Running
To start the HTTP server and accept requests to the API, run:
```
./push --config=/path/to/config.conf
```
By default, `default.conf` file in the same directory as the binary will be used as configuration.

## Tests
There are unit tests for `service` and `storage` packages.

To run tests, switch to package's directory and run
```
go test -v --config=/path/to/config.conf
```

## Develop

1. Clone repository
```
git clone git@gitlab.fixb.com:mir/push.git
```

2. Run server
```
go run main.go --config=/path/to/config.conf
```
