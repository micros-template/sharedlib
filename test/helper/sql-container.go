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
	Context                                                   context.Context
	SharedNetwork, ImageName, ContainerName                   string
	SQLInitScriptPath, SQLInitInsideScriptPath, WaitingSignal string
	Env                                                       map[string]string
}

func StartSQLContainer(opt SQLParameterOption) (*SQLContainer, error) {

	initSqlContent, err := os.ReadFile(opt.SQLInitScriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read init SQL file: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Name:     opt.ContainerName,
		Image:    opt.ImageName,
		Env:      opt.Env,
		Networks: []string{opt.SharedNetwork},
		WaitingFor: wait.ForLog(opt.WaitingSignal).
			WithOccurrence(2).WithStartupTimeout(5 * time.Second),
		Files: []testcontainers.ContainerFile{
			{
				Reader:            strings.NewReader(string(initSqlContent)),
				ContainerFilePath: opt.SQLInitInsideScriptPath,
				FileMode:          0644,
			},
		},
		ExposedPorts: []string{},
	}

	container, err := testcontainers.GenericContainer(opt.Context, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start SQL container: %w", err)
	}

	_, err = container.Host(opt.Context)
	if err != nil {
		container.Terminate(opt.Context)
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
