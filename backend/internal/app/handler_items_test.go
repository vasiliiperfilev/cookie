package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"golang.org/x/exp/slices"
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
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode("")
		request, err := http.NewRequest(http.MethodPost, "/v1/items", requestBody)
		request.Header.Set("Authorization", "Bearer "+strings.Repeat(strconv.FormatInt(2, 10), 26))
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		tester.AssertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("can't POST item if not supplier", func(t *testing.T) {
		// not an actual supplier
		supplierId := int64(1)
		dto := data.PostItemDto{
			Unit:     "l",
			Size:     1,
			Name:     "milk",
			ImageUrl: "test",
		}
		request := createPostItemRequest(t, dto, supplierId)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		tester.AssertStatus(t, response.Code, http.StatusForbidden)
	})

	t.Run("can't POST item with empty name, empty unit, size < 0", func(t *testing.T) {
		supplierId := int64(2)
		dto := data.PostItemDto{
			Unit:     "",
			Size:     -1,
			Name:     "",
			ImageUrl: "",
		}
		want := []string{
			"unit", "size", "name", "imageUrl",
		}
		request := createPostItemRequest(t, dto, supplierId)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		tester.AssertStatus(t, response.Code, http.StatusUnprocessableEntity)
		assertContentType(t, response, app.JsonContentType)
		var errors app.ErrorResponse
		json.NewDecoder(response.Body).Decode(&errors)
		for k := range errors.Errors {
			if !slices.Contains(want, k) {
				t.Fatalf("Want %v error key but not found", k)
			}
		}
	})
}

func TestItemGet(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	item1 := data.Item{Id: 1, SupplierId: 2, Name: "Milk", Unit: "l", Size: 1, ImageUrl: "test"}
	itemModel := data.NewStubItemModel([]data.Item{
		item1,
	})
	models := data.Models{User: data.NewStubUserModel(generateUsers(4)), Item: itemModel}
	server := app.New(cfg, logger, models)

	t.Run("it GET item if exists", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/items/%v", item1.Id), nil)
		request.Header.Set("Authorization", "Bearer "+strings.Repeat("1", 26))
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		tester.AssertStatus(t, response.Code, http.StatusOK)
		var got data.Item
		json.NewDecoder(response.Body).Decode(&got)
		want := item1
		tester.AssertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		tester.AssertValue(t, got, want, "Expected same item")
	})

	t.Run("it 404 if item doesn't exist", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/items/%v", 234), nil)
		request.Header.Set("Authorization", "Bearer "+strings.Repeat("1", 26))
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		tester.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("it 401 if incorrect token", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/items/%v", 234), nil)
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		tester.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("it GET all items of supplier_id", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/items?supplierId=%v", item1.SupplierId), nil)
		request.Header.Set("Authorization", "Bearer "+strings.Repeat("1", 26))
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		tester.AssertStatus(t, response.Code, http.StatusOK)
		var got []data.Item
		json.NewDecoder(response.Body).Decode(&got)
		want := item1
		tester.AssertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		tester.AssertValue(t, got[0], want, "Expected same item")
	})
}

func TestItemGetAll(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	item1 := data.Item{Id: 1, SupplierId: 2, Name: "Milk", Unit: "l", Size: 1, ImageUrl: "test"}
	itemModel := data.NewStubItemModel([]data.Item{
		item1,
	})
	models := data.Models{User: data.NewStubUserModel(generateUsers(4)), Item: itemModel}
	server := app.New(cfg, logger, models)

	t.Run("it GET all items of supplier_id", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/items?supplierId=%v", item1.SupplierId), nil)
		request.Header.Set("Authorization", "Bearer "+strings.Repeat("1", 26))
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		tester.AssertStatus(t, response.Code, http.StatusOK)
		var got []data.Item
		json.NewDecoder(response.Body).Decode(&got)
		want := item1
		tester.AssertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		tester.AssertValue(t, got[0], want, "Expected same item")
	})

	t.Run("it return 404 if supplier doesn't exist", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/items?supplierId=%v", 234), nil)
		request.Header.Set("Authorization", "Bearer "+strings.Repeat("1", 26))
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		tester.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("it return 401 if not authed", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/items?supplierId=%v", item1.SupplierId), nil)
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		tester.AssertStatus(t, response.Code, http.StatusUnauthorized)
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