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
type SQLParameterOption struct {
	context                                                   context.Context
	sharedNetwork, imageName, containerName                   string
	sqlInitScriptPath, sqlInitInsideScriptPath, waitingSignal string
	env                                                       map[string]string
}

func StartSQLContainer(opt SQLParameterOption) (*SQLContainer, error) {

	initSqlContent, err := os.ReadFile(opt.sqlInitScriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read init SQL file: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Name:     opt.containerName,
		Image:    opt.imageName,
		Env:      opt.env,
		Networks: []string{opt.sharedNetwork},
		WaitingFor: wait.ForLog(opt.waitingSignal).
			WithOccurrence(2).WithStartupTimeout(5 * time.Second),
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(string(initSqlContent)),
				ContainerFilePath: opt.sqlInitInsideScriptPath,
				FileMode:          0644,
			},
		},
		ExposedPorts: []string{},
	}

	container, err := testcontainers.GenericContainer(opt.context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start SQL container: %w", err)
	}

	_, err = container.Host(opt.context)
	if err != nil {
		container.Terminate(opt.context)
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
