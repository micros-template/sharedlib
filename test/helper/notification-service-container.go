package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type NotificationServiceContainer struct {
	Container testcontainers.Container
}

type NotificationServiceParameterOption struct {
	context                                                context.Context
	sharedNetwork, imageName, containerName, waitingSignal string
	cmd                                                    []string
	env                                                    map[string]string
}

func StartNotificationServiceContainer(opt NotificationServiceParameterOption) (*NotificationServiceContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:       opt.containerName,
		Image:      opt.imageName,
		Env:        opt.env,
		Networks:   []string{opt.sharedNetwork},
		Cmd:        opt.cmd,
		WaitingFor: wait.ForLog(opt.waitingSignal).WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(opt.context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start mail container: %w", err)
	}

	return &NotificationServiceContainer{
		Container: container,
	}, nil
}

func (f *NotificationServiceContainer) Terminate(ctx context.Context) error {
	if f.Container != nil {
		return f.Container.Terminate(ctx)
	}
	return nil
}
