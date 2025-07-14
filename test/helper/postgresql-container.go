package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	Container testcontainers.Container
}

func StartPostgresContainer(ctx context.Context, sharedNetwork, name, version string) (*PostgresContainer, error) {
	image := fmt.Sprintf("postgres:%s", version)
	req := testcontainers.ContainerRequest{
		Name:         name,
		Image:        image,
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       viper.GetString("database.name"),
			"POSTGRES_USER":     viper.GetString("database.user"),
			"POSTGRES_PASSWORD": viper.GetString("database.password"),
		},
		Networks: []string{sharedNetwork},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).WithStartupTimeout(5 * time.Second),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      viper.GetString("script.init_sql"),
				ContainerFilePath: "/docker-entrypoint-initdb.d/init-db.sql",
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

	_, err = container.MappedPort(ctx, "5432")
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	return &PostgresContainer{
		Container: container,
	}, nil
}

func (p *PostgresContainer) Terminate(ctx context.Context) error {
	if p.Container != nil {
		return p.Container.Terminate(ctx)
	}
	return nil
}
