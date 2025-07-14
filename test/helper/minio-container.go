package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MinioContainer struct {
	testcontainers.Container
}

func StartMinioContainer(ctx context.Context, sharedNetwork, version string) (*MinioContainer, error) {
	minioImage := fmt.Sprintf("minio/minio:%s", version)
	req := testcontainers.ContainerRequest{
		Name:         "minio",
		Image:        minioImage,
		ExposedPorts: []string{"9000/tcp"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     viper.GetString("minio.credential.user"),
			"MINIO_ROOT_PASSWORD": viper.GetString("minio.credential.password"),
		},
		Networks:   []string{sharedNetwork},
		Cmd:        []string{"server", "/data"},
		WaitingFor: wait.ForLog("API:").WithStartupTimeout(30 * time.Second),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start minio container: %w", err)
	}

	_, err = container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}
	_, err = container.MappedPort(ctx, "9000")
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	return &MinioContainer{
		Container: container,
	}, nil
}

func (mc *MinioContainer) Terminate(ctx context.Context) error {
	if mc.Container != nil {
		return mc.Container.Terminate(ctx)
	}
	return nil
}
