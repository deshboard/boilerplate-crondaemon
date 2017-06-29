package app

import (
	"fmt"

	"github.com/go-kit/kit/log"
)

// Job is responsible for the main logic.
type Job struct {
	logger log.Logger
}

// NewJob returns a new Job
func NewJob(logger log.Logger) *Job {
	return &Job{logger}
}

// Run executes the main logic.
func (j *Job) Run() {
	fmt.Println("Hello, World!")
}
