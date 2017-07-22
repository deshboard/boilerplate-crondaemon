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
	"github.com/goph/stdlib/ext"
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
	logger := newLogger(config)
	defer ext.Close(logger)

	// Create a new error handler
	errorHandler := newErrorHandler(config, logger)
	defer ext.Close(errorHandler)

	// Register error handler to recover from panics
	defer emperror.HandleRecover(errorHandler)

	// Initiate metrics scope
	metrics := newMetrics(config)
	defer ext.Close(metrics)

	// Application context
	appCtx := &application{
		config:          config,
		logger:          logger,
		errorHandler:    errorHandler,
		healthCollector: healthz.Collector{},
		tracer:          newTracer(config),
		metrics:         metrics,
	}

	serverQueue := serverz.NewQueue(&serverz.Manager{Logger: logger})

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

	job := newJob(appCtx)
	defer ext.Close(job)

	if config.Debug {
		debugServer := newDebugServer(logger)
		serverQueue.Append(debugServer, config.DebugAddr)
		defer debugServer.Close()
	}

	healthServer, status := newHealthServer(appCtx)
	serverQueue.Prepend(healthServer, config.HealthAddr)
	defer healthServer.Close()

	errChan := serverQueue.Start()

	// Necessary for daemon mode
	quit := make(chan struct{})

	if false == config.Daemon {
		job.Run()
	} else {
		ticker := time.NewTicker(config.DaemonSchedule)

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
	}

	signalChan := make(chan os.Signal, 1)
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
