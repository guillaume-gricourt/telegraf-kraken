.PHONY: all
all: ticker depth spread
ticker:
	CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' -o bin/telegraf-kraken-ticker cmd/main-ticker.go

depth:
	CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' -o bin/telegraf-kraken-depth cmd/main-depth.go

spread:
	CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' -o bin/telegraf-kraken-spread cmd/main-spread.go

.PHONY: test
test:
	go test plugins/inputs/ticker/ticker_test.go plugins/inputs/ticker/ticker.go
	go test plugins/inputs/depth/depth_test.go plugins/inputs/depth/depth.go
	go test plugins/inputs/spread/spread_test.go plugins/inputs/spread/spread.go
