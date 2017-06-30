package app

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
)

// Job is responsible for the main logic.
type Job struct {
	logger       log.Logger
	errorHandler emperror.Handler
}

// NewJob returns a new Job
func NewJob(logger log.Logger, errorHandler emperror.Handler) *Job {
	return &Job{logger, errorHandler}
}

// Run executes the main logic.
func (j *Job) Run() {
	fmt.Println("Hello, World!")
}
