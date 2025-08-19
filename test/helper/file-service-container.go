package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type FileServiceParameterOption struct {
	Context                                                context.Context
	SharedNetwork, ImageName, ContainerName, WaitingSignal string
	Cmd                                                    []string
	Env                                                    map[string]string
}

type FileServiceContainer struct {
	Container testcontainers.Container
}

func StartFileServiceContainer(opt FileServiceParameterOption) (*FileServiceContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:         opt.ContainerName,
		Image:        opt.ImageName,
		Env:          opt.Env,
		Networks:     []string{opt.SharedNetwork},
		Cmd:          opt.Cmd,
		ExposedPorts: []string{},
		WaitingFor:   wait.ForLog(opt.WaitingSignal).WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(opt.Context, testcontainers.GenericContainerRequest{
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
