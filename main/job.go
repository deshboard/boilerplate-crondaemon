package main

import "github.com/deshboard/boilerplate-crondaemon/app"

// newJob returns a new Job.
func newJob(appCtx *application) *app.Job {
	job := app.NewJob()

	job.Logger = appCtx.logger
	job.ErrorHandler = appCtx.errorHandler

	return job
}
