package data_test

import (
	"fmt"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/database"
	"github.com/vasiliiperfilev/cookie/internal/tester"
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

func TestUserModel(t *testing.T) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		database.POSTGRES_USER,
		database.POSTGRES_PASSWORD,
		database.POSTGRES_PORT,
		database.POSTGRES_DB,
	)
	t.Run("it inserts and gets a user", func(t *testing.T) {
		db := database.PrepareTestDb(t, dsn)
		model := data.NewPsqlUserModel(db)
		database.ApplyFixtures(t, db, "../fixtures")
		insertedUser := data.User{
			Email:   "test@test.com",
			Type:    1,
			ImageId: "id",
		}
		insertedUser.Password.Set("pa5$wOrd123")
		err := model.Insert(&insertedUser)
		tester.AssertNoError(t, err)
		gotUser, err := model.GetByEmail(insertedUser.Email)
		tester.AssertNoError(t, err)
		assertUser(t, *gotUser, insertedUser)
	})

	t.Run("it inserts 2 users concurently", func(t *testing.T) {
		db := database.PrepareTestDb(t, dsn)
		model := data.NewPsqlUserModel(db)
		database.ApplyFixtures(t, db, "../fixtures")

		users := []data.User{
			{
				Email:   "test1@test.com",
				Type:    1,
				ImageId: "id",
			},
			{
				Email:   "test2@test.com",
				Type:    1,
				ImageId: "id",
			},
		}
		errs := make([]error, 2)
		errsChannel := make(chan error)

		for _, user := range users {
			go func(user data.User) {
				user.Password.Set("pa5$wOrd123")
				errsChannel <- model.Insert(&user)
			}(user)
		}

		for i := 0; i < len(errs); i++ {
			err := <-errsChannel
			tester.AssertNoError(t, err)
		}
	})

	t.Run("it updates user", func(t *testing.T) {
		db := database.PrepareTestDb(t, dsn)
		model := data.NewPsqlUserModel(db)
		database.ApplyFixtures(t, db, "../fixtures")
		insertedUser := data.User{
			Email:   "test@test.com",
			Type:    1,
			ImageId: "id",
		}
		insertedUser.Password.Set("pa5$wOrd123")
		//insert a user
		err := model.Insert(&insertedUser)
		tester.AssertNoError(t, err)
		// change user type
		newEmail := "newtest2@test.com"
		insertedUser.Email = newEmail
		err = model.Update(&insertedUser)
		tester.AssertNoError(t, err)
		// get updateduser from db
		gotUser, err := model.GetByEmail(insertedUser.Email)
		tester.AssertNoError(t, err)
		tester.AssertValue(t, gotUser.Email, newEmail, "Expect the updated user type")
		tester.AssertValue(t, gotUser.Email, insertedUser.Email, "Expect same emails")
		tester.AssertValue(t, gotUser.Id, insertedUser.Id, "Expect same id")
	})
}

func assertErrorKeys(t *testing.T, keys []string, errors map[string]string) {
	t.Helper()
	for _, v := range keys {
		if _, ok := errors[v]; !ok {
			t.Fatalf("Expected to have error in %v field", v)
		}
	}
}

func assertUser(t *testing.T, got data.User, want data.User) {
	tester.AssertValue(t, got.Email, want.Email, "Expect same emails")
	tester.AssertValue(t, got.Id, want.Id, "Expect same id")
	tester.AssertValue(t, got.Type, want.Type, "Expect same type")
	tester.AssertValue(t, got.CreatedAt, want.CreatedAt, "Expect same createdAt")
}