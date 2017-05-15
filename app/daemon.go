package app

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

// Job is responsible for the main logic.
type Job struct {
	logger logrus.FieldLogger
}

// NewJob returns a new Job
func NewJob(logger logrus.FieldLogger) *Job {
	return &Job{logger}
}

// Run executes the main logic.
func (j *Job) Run() {
	fmt.Println("Hello, World!")
}
