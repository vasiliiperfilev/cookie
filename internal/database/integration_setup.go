package database

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vasiliiperfilev/cookie/internal/migrate"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

const (
	POSTGRES_DB       = "cookie_test"
	POSTGRES_USER     = "cookie"
	POSTGRES_PASSWORD = "cookie"
	POSTGRES_PORT     = "54350"
)

func PrepareTestDb(t *testing.T, dsn string) *sql.DB {
	t.Helper()
	cfg := Config{
		MaxOpenConns: 25,
		MaxIdleConns: 25,
		MaxIdleTime:  "15m",
		Dsn:          dsn,
	}
	// start a container
	err := StartDockerPostgres(t)
	tester.AssertNoError(t, err)
	// open connection
	db, err := OpenDB(cfg)
	tester.AssertNoError(t, err)
	// migrations
	err = migrate.Up(cfg.Dsn)
	tester.AssertNoError(t, err)

	return db
}

// starts postgress container with default testing credentials
func StartDockerPostgres(t *testing.T) error {
	t.Helper()
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.8-alpine",
		ExposedPorts: []string{fmt.Sprintf("%s:%s", POSTGRES_PORT, "5432")},
		WaitingFor:   wait.ForExposedPort(),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      "../../scripts/db/init-db.sql",
				ContainerFilePath: "/docker-entrypoint-initdb.d/init-db.sql",
				FileMode:          0755,
			},
		},
		Env: map[string]string{
			"POSTGRES_DB":       POSTGRES_DB,
			"POSTGRES_USER":     POSTGRES_USER,
			"POSTGRES_PASSWORD": POSTGRES_PASSWORD,
		},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return err
	}

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})
	return nil
}

// Clears up the DB and loads fixtures from filepath
func ApplyFixtures(t *testing.T, db *sql.DB, fixturesPath string) {
	t.Helper()
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(fixturesPath),
	)
	if err != nil {
		t.Fatalf("Unable to load fixtures, check DB setup or fixtures path, %s", err)
	}
	err = fixtures.Load()
	if err != nil {
		t.Fatalf("Unable to apply fixtures, check DB schema or fixtures file, %s", err)
	}
}
