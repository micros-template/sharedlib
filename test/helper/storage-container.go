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
type StorageParameterOption struct {
	Context                                                context.Context
	SharedNetwork, ImageName, ContainerName, WaitingSignal string
	Cmd                                                    []string
	Env                                                    map[string]string
}

func StartStorageContainer(opt StorageParameterOption) (*StorageContainer, error) {
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
		return nil, fmt.Errorf("failed to start storage container: %w", err)
	}

	_, err = container.Host(opt.Context)
	if err != nil {
		container.Terminate(opt.Context)
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
