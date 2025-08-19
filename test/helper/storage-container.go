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
	context                                                context.Context
	sharedNetwork, imageName, containerName, waitingSignal string
	cmd                                                    []string
	env                                                    map[string]string
}

func StartStorageContainer(opt StorageParameterOption) (*StorageContainer, error) {
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
		return nil, fmt.Errorf("failed to start storage container: %w", err)
	}

	_, err = container.Host(opt.context)
	if err != nil {
		container.Terminate(opt.context)
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
