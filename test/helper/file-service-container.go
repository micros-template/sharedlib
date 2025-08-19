package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type FileServiceParameterOption struct {
	context                                                context.Context
	sharedNetwork, imageName, containerName, waitingSignal string
	cmd                                                    []string
	env                                                    map[string]string
}

type FileServiceContainer struct {
	Container testcontainers.Container
}

func StartFileServiceContainer(opt FileServiceParameterOption) (*FileServiceContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:         opt.containerName,
		Image:        opt.imageName,
		Env:          opt.env,
		Networks:     []string{opt.sharedNetwork},
		Cmd:          opt.cmd,
		ExposedPorts: []string{},
		WaitingFor:   wait.ForLog(opt.waitingSignal).WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(opt.context, testcontainers.GenericContainerRequest{
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
