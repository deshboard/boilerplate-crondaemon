package main

import (
	"flag"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/sagikazarmark/healthz"
	"github.com/sagikazarmark/utilz/errors"
	"github.com/sagikazarmark/utilz/util"
)

// Global context variables
var (
	config           = &Configuration{}
	logger           = logrus.New().WithField("service", ServiceName)
	shutdownManager  = util.NewShutdownManager(errors.NewLogHandler(logger))
	checkerCollector = healthz.Collector{}
)

func init() {
	// Load configuration from environment
	err := envconfig.Process("", config)
	if err != nil {
		logger.Fatal(err)
	}

	defaultAddr := ""

	// Listen on loopback interface in development mode
	if config.Environment == "development" {
		defaultAddr = "127.0.0.1"
	}

	// Load flags into configuration
	flag.BoolVar(&config.Daemon, "daemon", false, "Start as daemon.")
	flag.StringVar(&config.HealthAddr, "health", defaultAddr+":10000", "Health service address.")
	flag.StringVar(&config.DebugAddr, "debug", defaultAddr+":10001", "Debug service address.")
	flag.DurationVar(&config.ShutdownTimeout, "shutdown", 2*time.Second, "Shutdown timeout.")
}
