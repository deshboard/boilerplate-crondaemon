package main

import (
	"github.com/deshboard/boilerplate-crondaemon/app"
	"github.com/go-kit/kit/log"
)

func newJob(logger log.Logger) *app.Job {
	return app.NewJob(logger)
}
