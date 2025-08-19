package helper

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MessageQueueContainer struct {
	Container testcontainers.Container
}
type MessageQueueParameterOption struct {
	Context                                         context.Context
	SharedNetwork, ImageName, ContainerName         string
	MQConfigPath, MQInsideConfigPath, WaitingSignal string
	MappedPort, Cmd                                 []string
	Env                                             map[string]string
}

func StartMessageQueueContainer(opt MessageQueueParameterOption) (*MessageQueueContainer, error) {

	natsConfigContent, err := os.ReadFile(opt.MQConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read NATS config file: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Name:         opt.ContainerName,
		Image:        opt.ImageName,
		ExposedPorts: opt.MappedPort,
		WaitingFor:   wait.ForLog(opt.WaitingSignal),
		Env:          opt.Env,
		Networks:     []string{opt.SharedNetwork},
		Cmd:          opt.Cmd,
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(string(natsConfigContent)),
				ContainerFilePath: opt.MQInsideConfigPath,
				FileMode:          0644,
			},
		},
	}

	container, err := testcontainers.GenericContainer(opt.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start message queue container: %w", err)
	}

	_, err = container.Host(opt.Context)
	if err != nil {
		container.Terminate(opt.Context)
		return nil, err
	}

	return &MessageQueueContainer{
		Container: container,
	}, nil
}

func (n *MessageQueueContainer) Terminate(ctx context.Context) error {
	if n.Container != nil {
		return n.Container.Terminate(ctx)
	}
	return nil
}
