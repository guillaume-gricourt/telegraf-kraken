# telegraf-kraken

This plugin will pull ticker stats from [Kraken API](https://docs.kraken.com/rest/#section/General-Usage/Support).

## Configuration

The minimal configuration expects the `pairs` to be set.

```toml
[[inputs.ticker]]
  ## Default is "https://api.kraken.com/0/public".
  # url = "https://api.kraken.com/0/public"

  ## Pairs to grab
  pairs = [""]

  ## HTTP method
  # method = "GET"

  ## Optional HTTP headers
  # headers = {"X-Special-Header" = "Special-Value"}

  ## Timeout for HTTP requests
  # timeout = "5s"
```

## Installation

* Clone the repository

```sh
git clone git@github.com:guillaume-gricourt/telegraf-kraken.git
```

* Build the "bin/telegraf-kraken-ticker-<label>" binary

The `label` refers to the endpoint of the API.
The labels available are: `ticker`

```sh
make <label>
```

* You should be able to call this from telegraf now using execd

```toml
[[inputs.execd]]
  command = ["/path/to/telegraf-kraken-<label>", "-poll_interval 1m"]
  signal = "none"
```

This self-contained plugin is based on the documentations of [Execd Go Shim](https://github.com/influxdata/telegraf/blob/effe112473a6bd8991ef8c12e293353c92f1d538/plugins/common/shim/README.md)
