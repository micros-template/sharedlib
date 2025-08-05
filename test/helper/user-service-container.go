package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type UserServiceContainer struct {
	Container testcontainers.Container
}

func StartUserServiceContainer(ctx context.Context, sharedNetwork, version string) (*UserServiceContainer, error) {
	image := fmt.Sprintf("10.1.20.130:5001/dropping/user-service:%s", version)
	req := testcontainers.ContainerRequest{
		Name:         "test_user_service",
		Image:        image,
		Env:          map[string]string{"ENV": "test"},
		Networks:     []string{sharedNetwork},
		Cmd:          []string{"/user_service"},
		ExposedPorts: []string{},
		WaitingFor:   wait.ForLog("gRPC server running in port").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &UserServiceContainer{
		Container: container,
	}, nil
}

func (u *UserServiceContainer) Terminate(ctx context.Context) error {
	if u.Container != nil {
		return u.Container.Terminate(ctx)
	}
	return nil
}
