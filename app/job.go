package app

import (
	"fmt"

	"github.com/go-kit/kit/log"
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
	fmt.Println("Hello, World!")
}
