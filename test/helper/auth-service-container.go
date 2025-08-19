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

func StartAuthServiceContainer(ctx context.Context, sharedNetwork, imageName, containerName, waitingSignal string, cmd []string, env map[string]string) (*AuthServiceContainer, error) {
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
