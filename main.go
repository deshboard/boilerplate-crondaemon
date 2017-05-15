package main // import "github.com/deshboard/boilerplate-crondaemon"

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/deshboard/boilerplate-crondaemon/app"
	"github.com/sagikazarmark/healthz"
	"github.com/sagikazarmark/serverz"
)

func main() {
	defer shutdown.Handle()

	flag.Parse()

	if config.Daemon && config.DaemonSchedule <= 0 {
		logger.Fatal("Daemon mode requires the DAEMON_SCHEDULE environment variable to be set.")
	}

	mode := "cron"
	if config.Daemon {
		mode = "daemon"
	}

	logger.WithFields(logrus.Fields{
		"version":     app.Version,
		"commitHash":  app.CommitHash,
		"buildDate":   app.BuildDate,
		"environment": config.Environment,
		"mode":        mode,
	}).Printf("Starting %s", app.FriendlyServiceName)

	job := app.NewJob(logger)

	if false == config.Daemon {
		job.Run()
	} else {
		w := logger.Logger.WriterLevel(logrus.ErrorLevel)
		shutdown.Register(w.Close)

		serverManager := serverz.NewServerManager(logger)
		errChan := make(chan error, 10)
		signalChan := make(chan os.Signal, 1)

		var debugServer serverz.Server
		if config.Debug {
			debugServer = &serverz.NamedServer{
				Server: &http.Server{
					Handler:  http.DefaultServeMux,
					ErrorLog: log.New(w, "debug: ", 0),
				},
				Name: "debug",
			}
			shutdown.RegisterAsFirst(debugServer.Close)

			go serverManager.ListenAndStartServer(debugServer, config.DebugAddr)(errChan)
		}

		status := healthz.NewStatusChecker(healthz.Healthy)
		readiness := status
		healthHandler := healthz.NewHealthServiceHandler(healthz.NewCheckers(), readiness)
		healthServer := &serverz.NamedServer{
			Server: &http.Server{
				Handler:  healthHandler,
				ErrorLog: log.New(w, "health: ", 0),
			},
			Name: "health",
		}
		shutdown.RegisterAsFirst(healthServer.Close)

		go serverManager.ListenAndStartServer(healthServer, config.HealthAddr)(errChan)

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

				if err != nil {
					logger.Error(err)
				} else {
					logger.Warning("Error channel received non-error value")
				}

				// Break the loop, proceed with regular shutdown
				break MainLoop
			case s := <-signalChan:
				logger.Infof(fmt.Sprintf("Captured %v", s))
				status.SetStatus(healthz.Unhealthy)
				close(quit)

				ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
				wg := &sync.WaitGroup{}

				if config.Debug {
					go serverManager.StopServer(debugServer, wg)(ctx)
				}
				go serverManager.StopServer(healthServer, wg)(ctx)

				wg.Wait()

				// Cancel context if shutdown completed earlier
				cancel()

				// Break the loop, proceed with regular shutdown
				break MainLoop
			}
		}
	}
}
