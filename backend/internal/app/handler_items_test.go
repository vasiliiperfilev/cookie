package app_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestItemPost(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	itemModel := data.NewStubItemModel([]data.Item{})
	models := data.Models{User: data.NewStubUserModel(generateUsers(4)), Item: itemModel}
	server := app.New(cfg, logger, models)

	t.Run("it POST item with correct values", func(t *testing.T) {
		itemId := int64(1)
		itemInput := data.PostItemDto{
			SupplierId: 2,
			Unit:       "l",
			Size:       1,
			Name:       "milk",
			ImageUrl:   "test",
		}
		want := data.Item{
			Id:         itemId,
			SupplierId: itemInput.SupplierId,
			Unit:       itemInput.Unit,
			Size:       itemInput.Size,
			Name:       itemInput.Name,
			ImageUrl:   itemInput.ImageUrl,
		}
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(itemInput)
		request, err := http.NewRequest(http.MethodPost, "/v1/items", requestBody)
		request.Header.Set("Authorization", "Bearer "+strings.Repeat("2", 26))
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		tester.AssertStatus(t, response.Code, http.StatusCreated)
		assertContentType(t, response, app.JsonContentType)
		var got data.Item
		json.NewDecoder(response.Body).Decode(&got)
		if got != want {
			t.Fatalf("Want %v, got %v", want, got)
		}
		got, err = itemModel.GetById(itemId)
		tester.AssertNoError(t, err)
		if got != want {
			t.Fatalf("Want %v, got %v", want, got)
		}
	})

	t.Run("can't POST item with empty body", func(t *testing.T) {

	})

	t.Run("can't POST item for another supplier", func(t *testing.T) {

	})

	t.Run("can't POST item if not supplier", func(t *testing.T) {

	})
}
