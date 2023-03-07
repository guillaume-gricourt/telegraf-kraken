.PHONY: all
all: ticker depth
ticker:
	CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' -o bin/telegraf-kraken-ticker cmd/main-ticker.go

depth:
	CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' -o bin/telegraf-kraken-depth cmd/main-depth.go

.PHONY: test
test:
	go test plugins/inputs/ticker/ticker_test.go plugins/inputs/ticker/ticker.go
	go test plugins/inputs/depth/depth_test.go plugins/inputs/depth/depth.go
