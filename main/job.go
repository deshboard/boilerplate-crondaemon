package main

import (
	"github.com/deshboard/boilerplate-crondaemon/app"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
)

// newJob returns a new Job.
func newJob(config *configuration, logger log.Logger, errorHandler emperror.Handler) *app.Job {
	return app.NewJob(logger, errorHandler)
}
