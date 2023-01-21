# telegraf-kraken

This plugin will pull ticker stats from (Kraken API)[https://docs.kraken.com/rest/#section/General-Usage/Support].

### Configuration

The minimal configuration expects the `pairs` to be set.
```toml
[[inputs.kraken]]
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

### Installation

* Clone the repo

```
git clone git@github.com:guillaume-gricourt/telegraf-kraken.git
```
* Build the "kraken" binary

```
make
```
* You should be able to call this from telegraf now using execd
```
[[inputs.execd]]
  command = ["/path/to/kraken", "-poll_interval 1m"]
  signal = "none"
```
This self-contained plugin is based on the documentations of [Execd Go Shim](https://github.com/influxdata/telegraf/blob/effe112473a6bd8991ef8c12e293353c92f1d538/plugins/common/shim/README.md)
