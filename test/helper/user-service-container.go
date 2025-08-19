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
	Context                                                context.Context
	SharedNetwork, ImageName, ContainerName, WaitingSignal string
	Cmd                                                    []string
	Env                                                    map[string]string
}

func StartUserServiceContainer(opt UserServiceParameterOption) (*UserServiceContainer, error) {
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
