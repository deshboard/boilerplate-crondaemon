package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
	"github.com/goph/healthz"
	"github.com/goph/serverz"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	config := &configuration{}

	// Load configuration from environment
	err := envconfig.Process("", config)
	if err != nil {
		panic(err)
	}

	// Load configuration from flags
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config.flags(flags)
	flags.Parse(os.Args[1:])

	if config.Daemon && config.DaemonSchedule <= 0 {
		panic("Daemon mode requires the DAEMON_SCHEDULE environment variable to be set.")
	}

	// Create a new logger
	logger, closer := newLogger(config)
	defer closer.Close()

	// Create a new error handler
	errorHandler, closer := newErrorHandler(config, logger)
	defer closer.Close()

	// Register error handler to recover from panics
	defer emperror.HandleRecover(errorHandler)

	mode := "cron"
	if config.Daemon {
		mode = "daemon"
	}

	level.Info(logger).Log(
		"msg", fmt.Sprintf("Starting %s", FriendlyServiceName),
		"version", Version,
		"commitHash", CommitHash,
		"buildDate", BuildDate,
		"environment", config.Environment,
		"mode", mode,
	)

	metricsReporter := newMetricsReporter(config)
	job, closer := newJob(config, logger, errorHandler, metricsReporter)
	defer closer.Close()

	if false == config.Daemon {
		job.Run()
	} else {
		healthCollector := healthz.Collector{}

		serverQueue := serverz.NewQueue(&serverz.Manager{Logger: logger})
		signalChan := make(chan os.Signal, 1)

		if config.Debug {
			debugServer := newDebugServer(logger)
			serverQueue.Append(debugServer, config.DebugAddr)
			defer debugServer.Close()
		}

		healthServer, status := newHealthServer(logger, healthCollector, metricsReporter)
		serverQueue.Prepend(healthServer, config.HealthAddr)
		defer healthServer.Close()

		errChan := serverQueue.Start()

		ticker := time.NewTicker(config.DaemonSchedule)
		quit := make(chan struct{})

		go func() {
			for {
				select {
				case <-quit:
					return
				case <-ticker.C:
					go job.Run()
				}
			}
		}()

		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	MainLoop:
		for {
			select {
			case err := <-errChan:
				status.SetStatus(healthz.Unhealthy)
				close(quit)
				level.Debug(logger).Log("msg", "Error received from error channel")
				emperror.HandleIfErr(errorHandler, err)

				// Break the loop, proceed with regular shutdown
				break MainLoop
			case s := <-signalChan:
				level.Info(logger).Log("msg", fmt.Sprintf("Captured %v", s))
				status.SetStatus(healthz.Unhealthy)
				close(quit)

				level.Debug(logger).Log(
					"msg", "Shutting down with timeout",
					"timeout", config.ShutdownTimeout,
				)

				ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)

				err := serverQueue.Stop(ctx)
				if err != nil {
					errorHandler.Handle(err)
				}

				// Cancel context if shutdown completed earlier
				cancel()

				// Break the loop, proceed with regular shutdown
				break MainLoop
			}
		}
	}
}
