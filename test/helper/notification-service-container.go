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
	Context                                                context.Context
	SharedNetwork, ImageName, ContainerName, WaitingSignal string
	Cmd                                                    []string
	Env                                                    map[string]string
}

func StartNotificationServiceContainer(opt NotificationServiceParameterOption) (*NotificationServiceContainer, error) {
	req := testcontainers.ContainerRequest{
		Name:       opt.ContainerName,
		Image:      opt.ImageName,
		Env:        opt.Env,
		Networks:   []string{opt.SharedNetwork},
		Cmd:        opt.Cmd,
		WaitingFor: wait.ForLog(opt.WaitingSignal).WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(opt.Context, testcontainers.GenericContainerRequest{
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
