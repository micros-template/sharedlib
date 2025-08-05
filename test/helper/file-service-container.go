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

func StartFileServiceContainer(ctx context.Context, sharedNetwork, version string) (*FileServiceContainer, error) {
	image := fmt.Sprintf("10.1.20.130:5001/dropping/file-service:%s", version)
	req := testcontainers.ContainerRequest{
		Name:  "test_file_service",
		Image: image,
		Env:        map[string]string{"ENV": "test"},
		Networks:   []string{sharedNetwork},
		Cmd:        []string{"/file_service"},
		ExposedPorts: []string{},
		WaitingFor: wait.ForLog("gRPC server running in port").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
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
