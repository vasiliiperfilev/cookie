package app_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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
		supplierId := int64(2)
		dto := data.PostItemDto{
			Unit:     "l",
			Size:     1,
			Name:     "milk",
			ImageUrl: "test",
		}
		want := data.Item{
			Id:         itemId,
			SupplierId: supplierId,
			Unit:       dto.Unit,
			Size:       dto.Size,
			Name:       dto.Name,
			ImageUrl:   dto.ImageUrl,
		}
		request := createPostItemRequest(t, dto, supplierId)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		tester.AssertStatus(t, response.Code, http.StatusCreated)
		assertContentType(t, response, app.JsonContentType)
		assertItemResponse(t, response, want)
		asserItemInModel(t, itemModel, itemId, want)
	})

	t.Run("can't POST unathorized", func(t *testing.T) {
		dto := data.PostItemDto{
			Unit:     "l",
			Size:     1,
			Name:     "milk",
			ImageUrl: "test",
		}
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(dto)
		request, err := http.NewRequest(http.MethodPost, "/v1/items", requestBody)
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		tester.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("can't POST item with empty body", func(t *testing.T) {

	})

	t.Run("can't POST item for another supplier", func(t *testing.T) {

	})

	t.Run("can't POST item if not supplier", func(t *testing.T) {

	})

	t.Run("can't POST item with empty name, empty unit, size < 0", func(t *testing.T) {

	})
}

func asserItemInModel(t *testing.T, itemModel *data.StubItemModel, itemId int64, want data.Item) {
	got, err := itemModel.GetById(itemId)
	tester.AssertNoError(t, err)
	if got != want {
		t.Fatalf("Want %v, got %v", want, got)
	}
}

func assertItemResponse(t *testing.T, response *httptest.ResponseRecorder, want data.Item) {
	var got data.Item
	json.NewDecoder(response.Body).Decode(&got)
	if got != want {
		t.Fatalf("Want %v, got %v", want, got)
	}
}

func createPostItemRequest(t *testing.T, dto data.PostItemDto, supplierId int64) *http.Request {
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(dto)
	request, err := http.NewRequest(http.MethodPost, "/v1/items", requestBody)
	tester.AssertNoError(t, err)
	request.Header.Set("Authorization", "Bearer "+strings.Repeat(strconv.FormatInt(supplierId, 10), 26))
	return request
}
