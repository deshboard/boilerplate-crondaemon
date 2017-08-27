package app

import (
	"github.com/go-kit/kit/log/level"
	"github.com/goph/stdlib/errors"
	"github.com/goph/stdlib/log"
)

// JobOption sets options in the Job.
type JobOption func(j *Job)

// Logger returns a JobOption that sets the logger for the Job.
func Logger(l log.Logger) JobOption {
	return func(j *Job) {
		j.logger = l
	}
}

// ErrorHandler returns a JobOption that sets the error handler for the Job.
func ErrorHandler(l errors.Handler) JobOption {
	return func(j *Job) {
		j.errorHandler = l
	}
}

// Job is responsible for the main logic.
type Job struct {
	logger       log.Logger
	errorHandler errors.Handler
}

// NewJob returns a new Job
func NewJob(opts ...JobOption) *Job {
	j := new(Job)

	for _, opt := range opts {
		opt(j)
	}

	// Default logger
	if j.logger == nil {
		j.logger = log.NewNopLogger()
	}

	// Default error handler
	if j.errorHandler == nil {
		j.errorHandler = errors.NewNopHandler()
	}

	return j
}

// Run executes the main logic.
func (j *Job) Run() {
	level.Info(j.logger).Log("msg", "Hello, World!")
}
