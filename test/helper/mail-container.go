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

type MailParameterOption struct {
	context                                                context.Context
	sharedNetwork, imageName, containerName, waitingSignal string
	mappedPort                                             []string
}

func StartMailContainer(opt MailParameterOption) (*MailContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:         opt.containerName,
		Image:        opt.imageName,
		ExposedPorts: opt.mappedPort,
		Networks:     []string{opt.sharedNetwork},
		WaitingFor:   wait.ForListeningPort(nat.Port(opt.waitingSignal)),
	}
	container, err := testcontainers.GenericContainer(opt.context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start mail container: %w", err)
	}

	_, err = container.Host(opt.context)
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
