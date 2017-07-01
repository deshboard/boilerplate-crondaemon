package app

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/goph/emperror"
)

// Job is responsible for the main logic.
type Job struct {
	Logger       log.Logger
	ErrorHandler emperror.Handler
}

// NewJob returns a new Job
func NewJob() *Job {
	return &Job{
		Logger:       log.NewNopLogger(),
		ErrorHandler: emperror.NewNullHandler(),
	}
}

// Run executes the main logic.
func (j *Job) Run() {
	level.Info(j.Logger).Log("msg", "Hello, World!")
}
