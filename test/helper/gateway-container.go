package helper

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type GatewayContainer struct {
	Container testcontainers.Container
}

func StartGatewayContainer(ctx context.Context, sharedNetwork, version string) (*GatewayContainer, error) {
	image := fmt.Sprintf("nginx:%s", version)

	nginxConfigPath := viper.GetString("script.nginx")
	nginxConfigContent, err := os.ReadFile(nginxConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read nginx config file: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Name:         "test_gateway",
		Image:        image,
		ExposedPorts: []string{"9090:80/tcp"},
		WaitingFor:   wait.ForLog("Configuration complete; ready for start up"),
		Networks:     []string{sharedNetwork},
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(string(nginxConfigContent)),
				ContainerFilePath: "/etc/nginx/conf.d/default.conf",
				FileMode:          0644,
			},
		},
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
		container.Terminate(ctx)
		return nil, err
	}

	_, err = container.MappedPort(ctx, "4221")
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
