package app_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	testfixtures "github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/db"
	"github.com/vasiliiperfilev/cookie/internal/migrate"
)

const (
	POSTGRES_DB       = "cookie_test"
	POSTGRES_USER     = "cookie"
	POSTGRES_PASSWORD = "cookie"
	POSTGRES_PORT     = "54350"
)

// can register and log in
func TestIntegrationUserPost(t *testing.T) {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_PORT, POSTGRES_DB)
	db := prepareTestDb(t, dsn)
	server := prepareServer(db, 4000)

	t.Run("it allows registration with correct values", func(t *testing.T) {
		applyFixtures(t, db, "../fixtures")
		userInput := data.RegisterUserInput{
			Email:    "test@nowhere.com",
			Password: "test123!A",
			Type:     1,
			ImageId:  "imageid",
		}
		expectedResponse := data.User{
			Email:   userInput.Email,
			Type:    userInput.Type,
			ImageId: userInput.ImageId,
		}
		response := registerUser(t, server, userInput)
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		assertRegisterResponse(t, response.Body, expectedResponse)
	})
}

func prepareServer(db *sql.DB, port int) *app.Application {
	cfg := app.Config{Port: port, Env: "development"}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	models := data.NewModels(db)

	server := app.New(cfg, logger, models)
	return server
}

func registerUser(t *testing.T, server http.Handler, input data.RegisterUserInput) *httptest.ResponseRecorder {
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(input)

	request, err := http.NewRequest(http.MethodPost, "/v1/user", requestBody)
	assertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)

	return response
}

func prepareTestDb(t *testing.T, dsn string) *sql.DB {
	t.Helper()
	cfg := db.Config{
		MaxOpenConns: 25,
		MaxIdleConns: 25,
		MaxIdleTime:  "15m",
		Dsn:          dsn,
	}
	// start a container
	err := startDockerPostgres(t)
	assertNoError(t, err)
	// open connection
	db, err := db.OpenDB(cfg)
	assertNoError(t, err)
	// migrations
	err = migrate.Up(cfg.Dsn)
	assertNoError(t, err)

	return db
}

// starts postgress container with default testing credentials
func startDockerPostgres(t *testing.T) error {
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
func applyFixtures(t *testing.T, db *sql.DB, fixturesPath string) {
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
