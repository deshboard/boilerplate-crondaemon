package main

import "github.com/deshboard/boilerplate-crondaemon/app"

func newJob() *app.Job {
	return app.NewJob(logger)
}
