package main

import (
	"flag"
	"time"

	_ "expvar"
	_ "net/http/pprof"

	"github.com/Sirupsen/logrus"
	"github.com/deshboard/boilerplate-crondaemon/app"
	"github.com/evalphobia/logrus_fluent"
	"github.com/kelseyhightower/envconfig"
	"github.com/sagikazarmark/utilz/errors"
	"github.com/sagikazarmark/utilz/util"
	"gopkg.in/airbrake/gobrake.v2"
	logrus_airbrake "gopkg.in/gemnasium/logrus-airbrake-hook.v2"
)

// Global context variables
var (
	config          = &app.Configuration{}
	logger          = logrus.New().WithField("service", app.ServiceName)
	shutdownManager = util.NewShutdownManager(errors.NewLogHandler(logger))
)

func init() {
	// Register shutdown handler in logrus
	logrus.RegisterExitHandler(shutdownManager.Shutdown)

	// Load configuration from environment
	err := envconfig.Process("", config)
	if err != nil {
		logger.Fatal(err)
	}

	// Log debug level messages if debug mode is turned on
	if config.Debug {
		logger.Logger.Level = logrus.DebugLevel
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

	// Initialize Airbrake
	if config.AirbrakeEnabled {
		airbrakeHook := logrus_airbrake.NewHook(config.AirbrakeProjectID, config.AirbrakeAPIKey, config.Environment)
		airbrake := airbrakeHook.Airbrake

		airbrake.SetHost(config.AirbrakeEndpoint)

		airbrake.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
			notice.Context["version"] = app.Version
			notice.Context["commit"] = app.CommitHash

			return notice
		})

		logger.Logger.Hooks.Add(airbrakeHook)
		shutdownManager.Register(airbrake.Close)
	}

	// Initialize Fluentd
	if config.FluentdEnabled {
		fluentdHook, err := logrus_fluent.New(config.FluentdHost, config.FluentdPort)
		if err != nil {
			logger.Panic(err)
		}

		// Configure fluent tag
		if app.FluentdTag != "" {
			fluentdHook.SetTag(app.FluentdTag)
		} else {
			fluentdHook.SetTag(app.ServiceName)
		}

		fluentdHook.AddFilter("error", logrus_fluent.FilterError)

		logger.Logger.Hooks.Add(fluentdHook)
		shutdownManager.Register(fluentdHook.Fluent.Close)
	}
}
