package main

import (
	"time"

	. "github.com/deshboard/boilerplate-crondaemon/app"
	"github.com/goph/serverz"
)

// newDaemonServer creates a new daemon server.
func newDaemonServer(app *application) serverz.Server {
	var ticker *time.Ticker

	if app.config.Daemon {
		ticker = time.NewTicker(app.config.DaemonSchedule)
	}

	return &serverz.AppServer{
		Server: &serverz.DaemonServer{
			Daemon: &serverz.CronDaemon{
				Job: NewService(
					Logger(app.logger),
					ErrorHandler(app.errorHandler),
				),
				Ticker: ticker,
			},
		},
		Name:   "daemon",
		Logger: app.logger,
	}
}
