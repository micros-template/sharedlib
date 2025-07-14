package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type AuthServiceContainer struct {
	Container testcontainers.Container
}

func StartAuthServiceContainer(ctx context.Context, sharedNetwork, version string) (*AuthServiceContainer, error) {
	image := fmt.Sprintf("auth_service:%s", version)
	req := testcontainers.ContainerRequest{
		Name:         "auth_service",
		Image:        image,
		ExposedPorts: []string{"8081:8081/tcp"},
		Env:          map[string]string{"ENV": "test"},
		Networks:     []string{sharedNetwork},
		Cmd:          []string{"/auth_service"},
		WaitingFor:   wait.ForLog("HTTP Server Starting in port").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &AuthServiceContainer{
		Container: container,
	}, nil
}

func (a *AuthServiceContainer) Terminate(ctx context.Context) error {
	if a.Container != nil {
		return a.Container.Terminate(ctx)
	}
	return nil
}
