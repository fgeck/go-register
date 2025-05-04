package integrationtest

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresConfig struct {
	Image    string
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

const (
	POSTGRES_DEFAULT_PORT    = "5432"
	POSTGRES_READY_LOG       = "database system is ready to accept connections"
	POSTGRES_READY_LOG_TIMES = 2
	POSTGRES_STARTUP_TIMEOUT = 5 * time.Second
)

// StartPostgres starts a PostgreSQL container for integration testing.
// It returns the container, the host, and the port of the PostgreSQL instance.
// The caller is responsible for terminating the container after use.
// The function also sets up the necessary environment variables for the database connection.
func StartPostgres(postgresCfg PostgresConfig) (testcontainers.Container, string, nat.Port, error) {
	ctx := context.Background()

	initScripts, err := getInitScripts()
	if err != nil {
		log.Printf("failed to get init scripts: %s", err)
		return nil, "", "", err
	}

	_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
	defer os.Unsetenv("TESTCONTAINERS_RYUK_DISABLED")
	postgresContainer, err := postgres.Run(ctx,
		postgresCfg.Image,
		postgres.WithInitScripts(initScripts...),
		postgres.WithDatabase(postgresCfg.Database),
		postgres.WithUsername(postgresCfg.Username),
		postgres.WithPassword(postgresCfg.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog(POSTGRES_READY_LOG).
				WithOccurrence(POSTGRES_READY_LOG_TIMES).
				WithStartupTimeout(POSTGRES_STARTUP_TIMEOUT)),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, "", "", err
	}
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Printf("failed to get container host: %s", err)
		return nil, "", "", err
	}
	port, err := postgresContainer.MappedPort(ctx, POSTGRES_DEFAULT_PORT)
	if err != nil {
		log.Printf("failed to get container port: %s", err)
		return nil, "", "", err
	}

	return postgresContainer, host, port, nil
}

func getInitScripts() ([]string, error) {
	var initScripts []string
	err := filepath.Walk("../migrations/postgresql", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		initScripts = append(initScripts, path)
		return nil
	})

	return initScripts, err
}
