package helper

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type NatsContainer struct {
	Container testcontainers.Container
}

func StartMessageQueueContainer(ctx context.Context, sharedNetwork, imageName, containerName, waitingSignal, mqConfigPath, mqInsideConfigPath string, mappedPort, cmd []string, env map[string]string) (*NatsContainer, error) {

	natsConfigContent, err := os.ReadFile(mqConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read NATS config file: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Name:         containerName,
		Image:        imageName,
		ExposedPorts: mappedPort,
		WaitingFor:   wait.ForLog(waitingSignal),
		Env:          env,
		Networks:     []string{sharedNetwork},
		Cmd:          cmd,
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(string(natsConfigContent)),
				ContainerFilePath: mqInsideConfigPath,
				FileMode:          0644,
			},
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start message queue container: %w", err)
	}

	_, err = container.Host(ctx)
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
