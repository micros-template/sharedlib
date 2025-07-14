package helper

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type NatsContainer struct {
	Container testcontainers.Container
}

func StartNatsContainer(ctx context.Context, sharedNetwork, version string) (*NatsContainer, error) {
	image := fmt.Sprintf("nats:%s", version)
	req := testcontainers.ContainerRequest{
		Name:         "nats",
		Image:        image,
		ExposedPorts: []string{"4221/tcp"},
		WaitingFor:   wait.ForLog("Server is ready"),
		Env: map[string]string{
			"NATS_USER":     viper.GetString("minio.credential.user"),
			"NATS_PASSWORD": viper.GetString("minio.credential.password"),
		},
		Networks: []string{sharedNetwork},
		Cmd: []string{
			"-c", "/etc/nats/nats.conf",
			"--name", "nats",
			"-p", "4221",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      viper.GetString("script.nats_server"),
				ContainerFilePath: "/etc/nats/nats.conf",
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

	return &NatsContainer{
		Container: container,
	}, nil
}

func (n *NatsContainer) Terminate(ctx context.Context) error {
	if n.Container != nil {
		return n.Container.Terminate(ctx)
	}
	return nil
}
