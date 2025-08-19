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
type CacheParameterOption struct {
	context                                                context.Context
	sharedNetwork, imageName, containerName, waitingSignal string
	cmd                                                    []string
	env                                                    map[string]string
}

func StartCacheContainer(opt CacheParameterOption) (*CacheContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:         opt.containerName,
		Image:        opt.imageName,
		Env:          opt.env,
		Networks:     []string{opt.sharedNetwork},
		Cmd:          opt.cmd,
		WaitingFor:   wait.ForListeningPort(nat.Port(opt.waitingSignal)),
		ExposedPorts: []string{},
	}
	container, err := testcontainers.GenericContainer(opt.context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start cache container: %w", err)
	}

	_, err = container.Host(opt.context)
	if err != nil {
		container.Terminate(opt.context)
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
