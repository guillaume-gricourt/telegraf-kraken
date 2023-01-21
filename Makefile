.PHONY: all
ticker:
	go build -o bin/telegraf-kraken-ticker cmd/main-ticker.go

.PHONY: test
test:
	go test plugins/inputs/ticker/ticker_test.go plugins/inputs/ticker/ticker.go
