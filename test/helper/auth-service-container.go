package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type AuthServiceContainer struct {
	Container testcontainers.Container
}
type AuthServiceParameterOption struct {
	Context                                                context.Context
	SharedNetwork, ImageName, ContainerName, WaitingSignal string
	Cmd                                                    []string
	Env                                                    map[string]string
}

func StartAuthServiceContainer(opt AuthServiceParameterOption) (*AuthServiceContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:       opt.ContainerName,
		Image:      opt.ImageName,
		Env:        opt.Env,
		Networks:   []string{opt.SharedNetwork},
		Cmd:        opt.Cmd,
		WaitingFor: wait.ForLog(opt.WaitingSignal).WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(opt.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start auth service container: %w", err)
	}

	return &AuthServiceContainer{
		Container: container,
	}, nil
}

func (a *AuthServiceContainer) Terminate(ctx context.Context) error {
	if a.Container != nil {
		return a.Container.Terminate(ctx)
	}
	return nil
}
