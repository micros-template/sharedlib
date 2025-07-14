package helper

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RedisContainer struct {
	testcontainers.Container
}

func StartRedisContainer(ctx context.Context, sharedNetwork, version string) (*RedisContainer, error) {
	image := fmt.Sprintf("redis:%s", version)
	req := testcontainers.ContainerRequest{
		Name:         "redis",
		Image:        image,
		ExposedPorts: []string{"6379/tcp"},
		Env: map[string]string{
			"REDIS_PASSWORD": viper.GetString("redis.password"),
		},
		Networks:   []string{sharedNetwork},
		Cmd:        []string{"redis-server", "--requirepass", viper.GetString("redis.password")},
		WaitingFor: wait.ForListeningPort("6379/tcp"),
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
	_, err = container.MappedPort(ctx, "6379")
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	return &RedisContainer{Container: container}, nil
}

func (r *RedisContainer) Terminate(ctx context.Context) error {
	if r.Container != nil {
		return r.Container.Terminate(ctx)
	}
	return nil
}
