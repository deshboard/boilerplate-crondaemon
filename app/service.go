package app

import (
	"github.com/deshboard/boilerplate-service/pkg/app"
	"github.com/go-kit/kit/log"
	"github.com/goph/emperror"
	"go.uber.org/dig"
)

// ServiceParams provides a set of dependencies for the service constructor.
type ServiceParams struct {
	dig.In

	Logger       log.Logger       `optional:"true"`
	ErrorHandler emperror.Handler `optional:"true"`
}

// NewService returns a new service instance.
func NewService(params ServiceParams) *app.Service {
	return app.NewService(
		app.Logger(params.Logger),
		app.ErrorHandler(params.ErrorHandler),
	)
}
