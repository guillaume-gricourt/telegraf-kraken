.PHONY: all
all: ticker
ticker:
	CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' -o bin/telegraf-kraken-ticker cmd/main-ticker.go

.PHONY: test
test:
	go test plugins/inputs/ticker/ticker_test.go plugins/inputs/ticker/ticker.go
