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
	context                                                context.Context
	sharedNetwork, imageName, containerName, waitingSignal string
	cmd                                                    []string
	env                                                    map[string]string
}

func StartAuthServiceContainer(opt AuthServiceParameterOption) (*AuthServiceContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:       opt.containerName,
		Image:      opt.imageName,
		Env:        opt.env,
		Networks:   []string{opt.sharedNetwork},
		Cmd:        opt.cmd,
		WaitingFor: wait.ForLog(opt.waitingSignal).WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(opt.context, testcontainers.GenericContainerRequest{
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
