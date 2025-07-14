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
	image := fmt.Sprintf("file_service:%s", version)
	req := testcontainers.ContainerRequest{
		Name:         "file_service",
		Image:        image,
		ExposedPorts: []string{"50052:50051/tcp"},
		Env:          map[string]string{"ENV": "test"},
		Networks:     []string{sharedNetwork},
		Cmd:          []string{"/file_service"},
		WaitingFor:   wait.ForLog("gRPC server running in port").WithStartupTimeout(30 * time.Second),
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
