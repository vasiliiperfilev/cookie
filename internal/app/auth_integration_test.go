package app_test

import (
	"context"
	"strings"
	"testing"

	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestRegistration(t *testing.T) {
	compose, err := tc.NewDockerCompose("../../docker-compose.yml")
	assertNoError(t, err)
	t.Cleanup(func() {
		assertNoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal))
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	err = compose.Up(ctx, tc.Wait(true))
	if err != nil && !strings.Contains(err.Error(), "exited (0)") {
		t.Fatalf("Expected to have no errors during Compose.Up()")
	}
}
