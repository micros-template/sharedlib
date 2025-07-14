package helper

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MailhogContainer struct {
	testcontainers.Container
}

func StartMailhogContainer(ctx context.Context, networkName, version string) (*MailhogContainer, error) {
	image := fmt.Sprintf("mailhog/mailhog:%s", version)
	req := testcontainers.ContainerRequest{
		Name:         "mailhog",
		Image:        image,
		ExposedPorts: []string{"1025:1025/tcp", "8025:8025/tcp"},
		Networks:     []string{networkName},
		WaitingFor:   wait.ForListeningPort("1025/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	_, err = container.Host(ctx)
	if err != nil {
		return nil, err
	}

	_, err = container.MappedPort(ctx, "1025")
	if err != nil {
		return nil, err
	}

	return &MailhogContainer{
		Container: container,
	}, nil
}

func (m *MailhogContainer) Terminate(ctx context.Context) error {
	if m.Container != nil {
		return m.Container.Terminate(ctx)
	}
	return nil
}
