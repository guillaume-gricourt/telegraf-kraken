package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/guillaume-gricourt/telegraf-kraken/plugins/inputs/kraken"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/common/shim"
)

var pollInterval = flag.Duration("poll_interval", 1*time.Minute, "how often to send metrics")
var pollIntervalDisabled = flag.Bool("poll_interval_disabled", false, "how often to send metrics")
var configFile = flag.String("config", "", "path to the config file for this plugin")
var usage = flag.Bool("usage", false, "print sample configuration")

func printConfig(name string, p telegraf.PluginDescriber) {
	fmt.Printf("[[inputs.%s]]", name)

	config := p.SampleConfig()
	if config != "" {
		fmt.Printf(config)
	} else {
		fmt.Printf("\n  # no configuration\n")
	}
}

func main() {
	flag.Parse()

	if *usage {
		printConfig("cryptodl", &kraken.Kraken{})
		os.Exit(0)
	}

	if *pollIntervalDisabled {
		*pollInterval = shim.PollIntervalDisabled
	}

	shim := shim.New()
	err := shim.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Err loading input: %v\n", err)
		os.Exit(1)
	}

	if err := shim.Run(*pollInterval); err != nil {
		fmt.Fprintf(os.Stderr, "Err: %s\n", err)
		os.Exit(1)
	}
}
