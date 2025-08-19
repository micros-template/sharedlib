package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type UserServiceContainer struct {
	Container testcontainers.Container
}
type UserServiceParameterOption struct {
	context                                                context.Context
	sharedNetwork, imageName, containerName, waitingSignal string
	cmd                                                    []string
	env                                                    map[string]string
}

func StartUserServiceContainer(opt UserServiceParameterOption) (*UserServiceContainer, error) {
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
		return nil, fmt.Errorf("failed to start user service container: %w", err)
	}

	return &UserServiceContainer{
		Container: container,
	}, nil
}

func (u *UserServiceContainer) Terminate(ctx context.Context) error {
	if u.Container != nil {
		return u.Container.Terminate(ctx)
	}
	return nil
}
