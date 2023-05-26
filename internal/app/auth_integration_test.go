package app_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	migrations "github.com/vasiliiperfilev/cookie/migrate"
)

// can register
// can't register with the same email
// can register and log in
func TestContainers(t *testing.T) {
	err := StartDockerPostgres(t)
	assertNoError(t, err)

	env := "testing"
	cfg := app.Config{Port: 4000, Env: env}
	cfg.Db.MaxOpenConns = 25
	cfg.Db.MaxIdleConns = 25
	cfg.Db.MaxIdleTime = "15m"
	cfg.Db.Dsn = "postgres://cookie:cookie@localhost:54350/cookie?sslmode=disable"
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := app.OpenDB(cfg)
	assertNoError(t, err)
	err = migrations.MigrateUp(cfg.Db.Dsn)
	assertNoError(t, err)

	query := `
        INSERT INTO user_type (type_name) 
        VALUES ('supplier')`
	_, err = db.Query(query)
	assertNoError(t, err)
	models := data.NewModels(db)
	server := app.New(cfg, logger, models)

	t.Run("it returns correct response", func(t *testing.T) {
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
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(userInput)

		request, err := http.NewRequest(http.MethodPost, "/v1/auth/register", requestBody)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertNoError(t, err)
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		assertRegisterResponse(t, response.Body, expectedResponse)
	})
}

func StartDockerPostgres(t *testing.T) error {
	t.Helper()
	ctx := context.Background()
	// dirname, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Current directory: %v\n", dirname)
	// dir := path.Join(dirname, "../../scripts/db")
	// assertNoError(t, err)
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.8-alpine",
		ExposedPorts: []string{"54350:5432"},
		WaitingFor:   wait.ForExposedPort(),
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      "../../scripts/db/install-extensions.sql",            // a directory
				ContainerFilePath: "/docker-entrypoint-initdb.d/install-extensions.sql", // important! its parent already exists
				FileMode:          0755,
			},
		},
		// Mounts: testcontainers.Mounts(testcontainers.ContainerMount{
		// 	Source: testcontainers.GenericBindMountSource{
		// 		HostPath: dir,
		// 	},
		// 	Target: "/docker-entrypoint-initdb.d",
		// }),
		Env: map[string]string{
			"POSTGRES_DB":       "cookie",
			"POSTGRES_USER":     "cookie",
			"POSTGRES_PASSWORD": "cookie",
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
