package main

import "github.com/deshboard/boilerplate-crondaemon/app"

func bootstrap() *app.Job {
	return app.NewJob(logger)
}
