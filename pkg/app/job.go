package app

import (
	"github.com/go-kit/kit/log/level"
)

// Run executes the main logic.
func (j *Service) Run() error {
	level.Info(j.logger).Log("msg", "Hello, World!")

	return nil
}
