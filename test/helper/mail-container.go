package helper

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MailContainer struct {
	testcontainers.Container
}

func StartMailContainer(ctx context.Context, networkName, imageName, containerName, waitingSignal string, mappedPort []string) (*MailContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:         containerName,
		Image:        imageName,
		ExposedPorts: mappedPort,
		Networks:     []string{networkName},
		WaitingFor:   wait.ForListeningPort(nat.Port(waitingSignal)),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start mail container: %w", err)
	}

	_, err = container.Host(ctx)
	if err != nil {
		return nil, err
	}

	_, err = container.MappedPort(ctx, nat.Port(waitingSignal))
	if err != nil {
		return nil, err
	}

	return &MailContainer{
		Container: container,
	}, nil
}

func (m *MailContainer) Terminate(ctx context.Context) error {
	if m.Container != nil {
		return m.Container.Terminate(ctx)
	}
	return nil
}
