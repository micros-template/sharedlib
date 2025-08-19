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
	Context                                                context.Context
	SharedNetwork, ImageName, ContainerName, WaitingSignal string
	Cmd                                                    []string
	Env                                                    map[string]string
}

func StartCacheContainer(opt CacheParameterOption) (*CacheContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:         opt.ContainerName,
		Image:        opt.ImageName,
		Env:          opt.Env,
		Networks:     []string{opt.SharedNetwork},
		Cmd:          opt.Cmd,
		WaitingFor:   wait.ForListeningPort(nat.Port(opt.WaitingSignal)),
		ExposedPorts: []string{},
	}
	container, err := testcontainers.GenericContainer(opt.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start cache container: %w", err)
	}

	_, err = container.Host(opt.Context)
	if err != nil {
		container.Terminate(opt.Context)
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
