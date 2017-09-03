package main

import (
	"time"

	"github.com/deshboard/boilerplate-crondaemon/app"
	"github.com/goph/serverz"
)

// newDaemonServer creates a new daemon server.
func newDaemonServer(appCtx *application) serverz.Server {
	var ticker *time.Ticker

	if appCtx.config.Daemon {
		ticker = time.NewTicker(appCtx.config.DaemonSchedule)
	}

	return &serverz.AppServer{
		Server: &serverz.DaemonServer{
			Daemon: &serverz.CronDaemon{
				Job: app.NewService(
					app.Logger(appCtx.logger),
					app.ErrorHandler(appCtx.errorHandler),
				),
				Ticker: ticker,
			},
		},
		Name:   "daemon",
		Logger: appCtx.logger,
	}
}
