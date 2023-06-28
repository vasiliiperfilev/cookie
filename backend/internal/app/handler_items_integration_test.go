package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/database"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestIntegrationPostItems(t *testing.T) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		database.POSTGRES_USER,
		database.POSTGRES_PASSWORD,
		database.POSTGRES_PORT,
		database.POSTGRES_DB,
	)
	cfg := database.Config{
		MaxOpenConns: 25,
		MaxIdleConns: 25,
		MaxIdleTime:  "15m",
		Dsn:          dsn,
	}
	db, err := database.OpenDB(cfg)
	tester.AssertNoError(t, err)

	t.Run("it POST and GET an item", func(t *testing.T) {
		server := app.PrepareIntegrationTestServer(db, 4000)
		email := "testItems@nowhere.com"
		password := "test123!A"
		registerInput := data.PostUserDto{
			Email:    email,
			Password: password,
			Name:     "test",
			Type:     data.SupplierUserType,
			ImageId:  "imageid",
		}
		// register a supplier
		user := mustRegisterUser(t, server, registerInput)
		loginInput := map[string]string{
			"Email":    email,
			"Password": password,
		}
		userToken := mustLoginUser(t, server, loginInput)
		// post item
		dto := data.PostItemDto{
			Unit:     "l",
			Size:     1,
			Name:     "milk",
			ImageUrl: "test",
		}
		want := data.Item{
			SupplierId: user.Id,
			Unit:       dto.Unit,
			Size:       dto.Size,
			Name:       dto.Name,
			ImageUrl:   dto.ImageUrl,
		}
		// assert item
		got := postItem(t, server, userToken.Token.Plaintext, dto)
		want.Id = got.Id
		tester.AssertValue(t, got, want, "Expected to have same item after POST")
		// get array of supplier's items
		items := getItemsBySupplierId(t, server, userToken, user.Id)
		tester.AssertValue(t, items[0], want, "Expected to have same item after GET all")
		// get one item
		got = getItem(t, server, userToken.Token.Plaintext, want.Id)
		tester.AssertValue(t, got, want, "Expected to have same item after GET")
		// put item
		dto.Unit = "l"
		dto.Name = "Juice"
		dto.Size = 2
		dto.ImageUrl = "test 2"
		want.Unit = "l"
		want.Name = "Juice"
		want.Size = 2
		want.ImageUrl = "test 2"
		request := putItemRequest(t, userToken.Token.Plaintext, dto, want.Id)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		tester.AssertStatus(t, response.Code, http.StatusOK)
		got = decodeResponse[data.Item](t, response)
		tester.AssertValue(t, got, want, "Expected to have same item after PUT")
		got = getItem(t, server, userToken.Token.Plaintext, want.Id)
		tester.AssertValue(t, got, want, "Expected to GET same item after PUT")
	})
}

func postItem(t *testing.T, server http.Handler, token string, dto data.PostItemDto) data.Item {
	t.Helper()
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(dto)

	request, err := http.NewRequest(http.MethodPost, "/v1/items", requestBody)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	tester.AssertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	tester.AssertStatus(t, response.Code, http.StatusCreated)
	var item data.Item
	json.NewDecoder(response.Body).Decode(&item)
	return item
}

func getItemsBySupplierId(t *testing.T, server *app.Application, token app.UserToken, supplierId int64) []data.Item {
	t.Helper()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/items?supplierId=%v", supplierId), nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token.Token.Plaintext))
	tester.AssertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	tester.AssertStatus(t, response.Code, http.StatusOK)
	var items []data.Item
	json.NewDecoder(response.Body).Decode(&items)

	return items
}

func getItem(t *testing.T, server *app.Application, token string, itemId int64) data.Item {
	t.Helper()

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/items/%v", itemId), nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	tester.AssertNoError(t, err)
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	tester.AssertStatus(t, response.Code, http.StatusOK)
	var item data.Item
	json.NewDecoder(response.Body).Decode(&item)

	return item
}

func putItemRequest(t *testing.T, token string, dto data.PostItemDto, itemId int64) *http.Request {
	t.Helper()
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(dto)

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/v1/items/%v", itemId), requestBody)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	tester.AssertNoError(t, err)
	return request
}

func decodeResponse[T any](t *testing.T, response *httptest.ResponseRecorder) T {
	var item T
	json.NewDecoder(response.Body).Decode(&item)
	return item
}
