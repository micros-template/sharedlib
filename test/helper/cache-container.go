package helper

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type CacheContainer struct {
	testcontainers.Container
}

func StartCacheContainer(ctx context.Context, sharedNetwork, imageName, containerName, waitingSignal string, cmd []string, env map[string]string) (*CacheContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:         containerName,
		Image:        imageName,
		Env:          env,
		Networks:     []string{sharedNetwork},
		Cmd:          cmd,
		WaitingFor:   wait.ForListeningPort(nat.Port(waitingSignal)),
		ExposedPorts: []string{},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start cache container: %w", err)
	}

	_, err = container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	return &CacheContainer{Container: container}, nil
}

func (r *CacheContainer) Terminate(ctx context.Context) error {
	if r.Container != nil {
		return r.Container.Terminate(ctx)
	}
	return nil
}
