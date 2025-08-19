package helper

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type GatewayContainer struct {
	Container testcontainers.Container
}

func StartGatewayContainer(ctx context.Context, sharedNetwork, imageName, containerName, nginxConfigPath, nginxInsideConfigPath, grpcErrorConfigPath, grpcErrorInsideConfigPath, waitingSignal string, mappedPort []string) (*GatewayContainer, error) {
	nginxConfigContent, err := os.ReadFile(nginxConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read nginx config file: %w", err)
	}

	grpcErrorConfigContent, err := os.ReadFile(grpcErrorConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read error.grpc_conf file: %w", err)
	}
	req := testcontainers.ContainerRequest{
		Name:         containerName,
		Image:        imageName,
		ExposedPorts: mappedPort,
		WaitingFor:   wait.ForLog(waitingSignal),
		Networks:     []string{sharedNetwork},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(string(nginxConfigContent)),
				ContainerFilePath: nginxInsideConfigPath,
				FileMode:          0644,
			},
			{
				Reader:            strings.NewReader(string(grpcErrorConfigContent)),
				ContainerFilePath: grpcErrorInsideConfigPath,
				FileMode:          0644,
			},
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start gateway container: %w", err)

	}

	_, err = container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	return &GatewayContainer{
		Container: container,
	}, nil
}

func (g *GatewayContainer) Terminate(ctx context.Context) error {
	if g.Container != nil {
		return g.Container.Terminate(ctx)
	}
	return nil
}
