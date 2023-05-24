package data_test

import (
	"fmt"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
)

func TestValidateRegisterUserInput(t *testing.T) {
	inputs := []struct {
		Input data.RegisterUserInput
		Keys  []string
	}{
		{Input: data.RegisterUserInput{
			Email:    "test@test.com",
			Password: "pa5$wOrd123",
			Type:     1,
			ImageId:  "testid",
		}, Keys: make([]string, 0)},
		{Input: data.RegisterUserInput{
			Email:    "test@test.com",
			Password: "",
			Type:     3,
			ImageId:  "",
		}, Keys: []string{"password", "type", "imageId"}},
		{Input: data.RegisterUserInput{
			Email:    "test-test.com",
			Password: "pa5swOrd123",
			Type:     1,
			ImageId:  "testId",
		}, Keys: []string{"email", "password"}},
		{Input: data.RegisterUserInput{
			Email:    "test@test.com",
			Password: "pa5$word123",
			Type:     1,
			ImageId:  "testId",
		}, Keys: []string{"password"}},
		{Input: data.RegisterUserInput{
			Email:    "test@test.com",
			Password: "pas$wOrdabc",
			Type:     1,
			ImageId:  "testId",
		}, Keys: []string{"password"}},
		{Input: data.RegisterUserInput{
			Email:    "test@test.com",
			Password: "PAS$WORD123",
			Type:     1,
			ImageId:  "testId",
		}, Keys: []string{"password"}},
	}

	for _, input := range inputs {
		t.Run(fmt.Sprintf("Incorrect keys %v", input.Keys), func(t *testing.T) {
			v := validator.New()
			data.ValidateRegisterUserInput(v, &input.Input)
			assertErrorKeys(t, input.Keys, v.Errors)
		})
	}
}

func assertErrorKeys(t *testing.T, keys []string, errors map[string]string) {
	t.Helper()
	for _, v := range keys {
		if _, ok := errors[v]; !ok {
			t.Fatalf("Expected to have error in %v field", v)
		}
	}
}
