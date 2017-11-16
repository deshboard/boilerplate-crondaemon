package app

import (
	"time"

	"github.com/deshboard/boilerplate-crondaemon/pkg/app"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	"github.com/goph/fxt/daemon"
	"go.uber.org/dig"
)

// ServiceParams provides a set of dependencies for the service constructor.
type ServiceParams struct {
	dig.In

	Config       *Config
	Logger       log.Logger       `optional:"true"`
	ErrorHandler emperror.Handler `optional:"true"`
}

// NewService returns a new service instance.
func NewService(params ServiceParams) daemon.Daemon {
	service := app.NewService(
		app.Logger(params.Logger),
		app.ErrorHandler(params.ErrorHandler),
	)

	var ticker *time.Ticker

	if params.Config.Daemon {
		ticker = time.NewTicker(params.Config.DaemonSchedule)
	}

	return &daemon.CronDaemon{
		Job:    service,
		Ticker: ticker,
	}
}
