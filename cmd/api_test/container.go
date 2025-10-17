package main

import (
	"context"

	"github.com/DTSL/golang-libraries/di"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// GetContainer returns fully initialized container.
func getContainer(ctx context.Context, app di.Application) (*dig.Container, error) {
	// Create Dependency Injection (DI) container
	container, err := di.GetContainer(ctx, app)
	if err != nil {
		return nil, errors.Wrap(err, "get di container")
	}

	// Provide resources
	for _, provide := range []struct {
		Name     string
		Resource any
		Options  []dig.ProvideOption
	}{
		{
			Name:     "router",
			Resource: newRouter,
		},
		{
			Name:     "test http handler",
			Resource: newhttpTestHandler,
		},
	} {
		provideErr := container.Provide(provide.Resource, provide.Options...)
		if provideErr != nil {
			return nil, errors.Wrapf(provideErr, "provide %s", provide.Name)
		}
	}

	return container, nil
}
