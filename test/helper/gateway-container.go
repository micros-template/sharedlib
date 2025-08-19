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
type GatewayParameterOption struct {
	context                                                       context.Context
	sharedNetwork, imageName, containerName                       string
	nginxConfigPath, nginxInsideConfigPath                        string
	grpcErrorConfigPath, grpcErrorInsideConfigPath, waitingSignal string
	mappedPort                                                    []string
}

func StartGatewayContainer(opt GatewayParameterOption) (*GatewayContainer, error) {
	nginxConfigContent, err := os.ReadFile(opt.nginxConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read nginx config file: %w", err)
	}

	grpcErrorConfigContent, err := os.ReadFile(opt.grpcErrorConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read error.grpc_conf file: %w", err)
	}
	req := testcontainers.ContainerRequest{
		Name:         opt.containerName,
		Image:        opt.imageName,
		ExposedPorts: opt.mappedPort,
		WaitingFor:   wait.ForLog(opt.waitingSignal),
		Networks:     []string{opt.sharedNetwork},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(string(nginxConfigContent)),
				ContainerFilePath: opt.nginxInsideConfigPath,
				FileMode:          0644,
			},
			{
				Reader:            strings.NewReader(string(grpcErrorConfigContent)),
				ContainerFilePath: opt.grpcErrorInsideConfigPath,
				FileMode:          0644,
			},
		},
	}

	container, err := testcontainers.GenericContainer(opt.context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start gateway container: %w", err)

	}

	_, err = container.Host(opt.context)
	if err != nil {
		container.Terminate(opt.context)
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
