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
	Context                                                context.Context
	SharedNetwork, ImageName, ContainerName, WaitingSignal string
	MappedPort                                             []string
}

func StartMailContainer(opt MailParameterOption) (*MailContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:         opt.ContainerName,
		Image:        opt.ImageName,
		ExposedPorts: opt.MappedPort,
		Networks:     []string{opt.SharedNetwork},
		WaitingFor:   wait.ForListeningPort(nat.Port(opt.WaitingSignal)),
	}
	container, err := testcontainers.GenericContainer(opt.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start mail container: %w", err)
	}

	_, err = container.Host(opt.Context)
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
