package main

import "github.com/deshboard/boilerplate-crondaemon/app"

// newJob returns a new Job.
func newJob(appCtx *application) *app.Job {
	job := app.NewJob(
		app.Logger(appCtx.logger),
		app.ErrorHandler(appCtx.errorHandler),
	)

	return job
}
