package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type FileServiceContainer struct {
	Container testcontainers.Container
}

func StartFileServiceContainer(ctx context.Context, sharedNetwork, imageName, containerName, waitingSignal string, cmd []string, env map[string]string) (*FileServiceContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:         containerName,
		Image:        imageName,
		Env:          env,
		Networks:     []string{sharedNetwork},
		Cmd:          cmd,
		ExposedPorts: []string{},
		WaitingFor:   wait.ForLog("waitingSignal").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start file service container: %w", err)
	}

	return &FileServiceContainer{
		Container: container,
	}, nil
}

func (f *FileServiceContainer) Terminate(ctx context.Context) error {
	if f.Container != nil {
		return f.Container.Terminate(ctx)
	}
	return nil
}
