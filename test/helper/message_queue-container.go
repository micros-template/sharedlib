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
	context                                         context.Context
	sharedNetwork, imageName, containerName         string
	mqConfigPath, mqInsideConfigPath, waitingSignal string
	mappedPort, cmd                                 []string
	env                                             map[string]string
}

func StartMessageQueueContainer(opt MessageQueueParameterOption) (*MessageQueueContainer, error) {

	natsConfigContent, err := os.ReadFile(opt.mqConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read NATS config file: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Name:         opt.containerName,
		Image:        opt.imageName,
		ExposedPorts: opt.mappedPort,
		WaitingFor:   wait.ForLog(opt.waitingSignal),
		Env:          opt.env,
		Networks:     []string{opt.sharedNetwork},
		Cmd:          opt.cmd,
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(string(natsConfigContent)),
				ContainerFilePath: opt.mqInsideConfigPath,
				FileMode:          0644,
			},
		},
	}

	container, err := testcontainers.GenericContainer(opt.context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start message queue container: %w", err)
	}

	_, err = container.Host(opt.context)
	if err != nil {
		container.Terminate(opt.context)
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
