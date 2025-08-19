package helper

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type SQLContainer struct {
	Container testcontainers.Container
}

func StartSQLContainer(ctx context.Context, sharedNetwork, imageName, containerName, sqlInitScriptPath, sqlInitInsideScriptPath, waitingSignal string, env map[string]string) (*SQLContainer, error) {

	initSqlContent, err := os.ReadFile(sqlInitScriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read init SQL file: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Name:     containerName,
		Image:    imageName,
		Env:      env,
		Networks: []string{sharedNetwork},
		WaitingFor: wait.ForLog(waitingSignal).
			WithOccurrence(2).WithStartupTimeout(5 * time.Second),
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(string(initSqlContent)),
				ContainerFilePath: sqlInitInsideScriptPath,
				FileMode:          0644,
			},
		},
		ExposedPorts: []string{},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start SQL container: %w", err)
	}

	_, err = container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, err
	}

	return &SQLContainer{
		Container: container,
	}, nil
}

func (p *SQLContainer) Terminate(ctx context.Context) error {
	if p.Container != nil {
		return p.Container.Terminate(ctx)
	}
	return nil
}
