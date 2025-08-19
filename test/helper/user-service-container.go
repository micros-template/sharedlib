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

func StartUserServiceContainer(ctx context.Context, sharedNetwork, imageName, containerName, waitingSignal string, cmd []string, env map[string]string) (*UserServiceContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:       containerName,
		Image:      imageName,
		Env:        env,
		Networks:   []string{sharedNetwork},
		Cmd:        cmd,
		WaitingFor: wait.ForLog(waitingSignal).WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
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
