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

func StartNotificationServiceContainer(ctx context.Context, sharedNetwork, version string) (*NotificationServiceContainer, error) {
	image := fmt.Sprintf("notification_service:%s", version)
	req := testcontainers.ContainerRequest{
		Name:       "notification_service",
		Image:      image,
		Env:        map[string]string{"ENV": "test"},
		Networks:   []string{sharedNetwork},
		Cmd:        []string{"/notification_service"},
		WaitingFor: wait.ForLog("subscriber for notification is running").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
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
