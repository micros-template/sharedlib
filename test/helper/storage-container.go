package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type StorageContainer struct {
	testcontainers.Container
}

func StartStorageContainer(ctx context.Context, sharedNetwork, imageName, containerName, waitingSignal string, cmd []string, env map[string]string) (*StorageContainer, error) {
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
		return nil, fmt.Errorf("failed to start storage container: %w", err)
	}

	_, err = container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	return &StorageContainer{
		Container: container,
	}, nil
}

func (mc *StorageContainer) Terminate(ctx context.Context) error {
	if mc.Container != nil {
		return mc.Container.Terminate(ctx)
	}
	return nil
}
