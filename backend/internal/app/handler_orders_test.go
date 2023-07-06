package app_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/vasiliiperfilev/cookie/internal/app"
	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/tester"
)

func TestOrderPost(t *testing.T) {
	cfg := app.Config{Port: 4000, Env: "development"}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	itemModel := data.NewStubItemModel([]data.Item{
		{
			Id:         1,
			SupplierId: 2,
		},
		{
			Id:         2,
			SupplierId: 2,
		},
		{
			Id:         3,
			SupplierId: 4,
		},
	})
	conversationModel := data.NewStubConversationModel(generateConversation(4))
	messageModel := data.NewStubMessageModel(generateConversation(4), []data.Message{{Id: 1, ConversationId: 1, PrevMessageId: 0}})
	orderModel := data.NewStubOrderModel([]data.Order{}, itemModel, conversationModel, messageModel)
	models := data.Models{
		Conversation: data.NewStubConversationModel(generateConversation(4)),
		User:         data.NewStubUserModel(generateUsers(4)),
		Item:         itemModel,
		Message:      messageModel,
		Order:        orderModel,
	}
	repositories := data.Repositories{Order: data.NewStubOrderRepository(orderModel, messageModel)}
	server := app.New(cfg, logger, models, repositories)

	t.Run("it POST order with correct values", func(t *testing.T) {
		clientId := int64(1)
		dto := data.PostOrderDto{
			ConversationId: 1,
			ItemIds:        []int64{1, 2},
		}
		request := createPostOrderRequest(t, dto, clientId)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		messages, err := messageModel.GetAllByConversationId(1)
		tester.AssertNoError(t, err)
		want := data.Order{
			ItemIds:   []int64{1, 2},
			StateId:   data.OrderStateCreated,
			MessageId: int64(len(messages) - 1), // order attached to last message
		}

		tester.AssertStatus(t, response.Code, http.StatusCreated)
		assertContentType(t, response, app.JsonContentType)
		got := parseOrderResponse(t, response)
		want.Id = got.Id
		assertOrder(t, got, want)
		assertOrderInModel(t, orderModel, got.Id, want)
	})

	t.Run("it 422 if POST order with non existing items", func(t *testing.T) {
		clientId := int64(1)
		dto := data.PostOrderDto{
			ConversationId: 1,
			ItemIds:        []int64{4, 5},
		}
		request := createPostOrderRequest(t, dto, clientId)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		messages, err := messageModel.GetAllByConversationId(1)
		tester.AssertNoError(t, err)
		wantMsgCount := len(messages)
		orders, err := orderModel.GetAllByUserId(1)
		tester.AssertNoError(t, err)
		wantOrderCount := len(orders)

		orders, err = orderModel.GetAllByUserId(1)
		tester.AssertNoError(t, err)
		gotOrderCount := len(orders)
		tester.AssertValue(t, gotOrderCount, wantOrderCount, "Expected to not have new orders")
		messages, err = messageModel.GetAllByConversationId(1)
		tester.AssertNoError(t, err)
		gotMsgCount := len(messages)
		tester.AssertValue(t, gotMsgCount, wantMsgCount, "Expected to not have new messages")
		tester.AssertStatus(t, response.Code, http.StatusUnprocessableEntity)
		assertContentType(t, response, app.JsonContentType)
	})

	t.Run("it 422 if POST order with items of different suppliers", func(t *testing.T) {

	})

	t.Run("it 401 if POST order unathorized", func(t *testing.T) {

	})
}

// func TestOrderGet(t *testing.T) {
// 	cfg := app.Config{Port: 4000, Env: "development"}
// 	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
// 	itemModel := data.NewStubItemModel([]data.Item{})
// 	models := data.Models{User: data.NewStubUserModel(generateUsers(4)), Item: itemModel}
// 	server := app.New(cfg, logger, models)

// 	t.Run("it GET order", func(t *testing.T) {

// 	})

// 	t.Run("it 404 if GET non-existing order", func(t *testing.T) {

// 	})

// 	t.Run("it 401 if GET order unathorized", func(t *testing.T) {

// 	})

// 	t.Run("it 403 if GET not own order", func(t *testing.T) {

// 	})
// }

// func TestOrderGetAll(t *testing.T) {
// 	cfg := app.Config{Port: 4000, Env: "development"}
// 	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
// 	itemModel := data.NewStubItemModel([]data.Item{})
// 	models := data.Models{User: data.NewStubUserModel(generateUsers(4)), Item: itemModel}
// 	server := app.New(cfg, logger, models)

// 	t.Run("it GET all orders of own id", func(t *testing.T) {

// 	})

// 	t.Run("it 404 if GET all orders of non-existing", func(t *testing.T) {

// 	})

// 	t.Run("it 401 if GET all orders unathorized", func(t *testing.T) {

// 	})

// 	t.Run("it 403 if GET all orders of not owning user", func(t *testing.T) {

// 	})
// }

// func TestOrderPut(t *testing.T) {
// 	cfg := app.Config{Port: 4000, Env: "development"}
// 	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
// 	itemModel := data.NewStubItemModel([]data.Item{})
// 	models := data.Models{User: data.NewStubUserModel(generateUsers(4)), Item: itemModel}
// 	server := app.New(cfg, logger, models)

// 	// Order states with owner:
// 	// Created - 1 (client)
// 	// Accepted - 2 (supplier)
// 	// Declined - 3 (supplier)
// 	// Fulfilled - 4 (supplier)
// 	// Confirmed fulfillment - 5 (client)
// 	// Supplier changes - 6 (supplier)
// 	// Client changes - 7 (client)

// 	t.Run("it 201 if supplier PUT order to accepted", func(t *testing.T) {

// 	})

// 	t.Run("it 201 if supplier PUT order to fulfilled", func(t *testing.T) {
// 		// stop at this point and implement frontend
// 	})

// 	t.Run("it 201 if client PUT order to confirm fulfielment", func(t *testing.T) {

// 	})

// 	t.Run("it 201 if supplier PUT order to suggest changes before fulfielment", func(t *testing.T) {

// 	})

// 	t.Run("it 201 if client PUT order to suggest changes before fulfielment", func(t *testing.T) {

// 	})

// 	t.Run("it 201 if client PUT order to accept supplier changes", func(t *testing.T) {

// 	})

// 	t.Run("it 201 if client PUT order to accept client changes", func(t *testing.T) {

// 	})

// 	t.Run("it 201 if supplier PUT order to decline before accepting/suggesting changes", func(t *testing.T) {

// 	})

// 	t.Run("it 400 if supplier PUT order to decline after fulfilled", func(t *testing.T) {

// 	})

// 	t.Run("it 400 if supplier PUT order to accepted after own changes", func(t *testing.T) {

// 	})

// 	t.Run("it 400 if client PUT order to accepted after own changes", func(t *testing.T) {

// 	})

// 	t.Run("it 400 if supplier PUT order to canceled by client", func(t *testing.T) {

// 	})

// 	t.Run("it 400 if client PUT order to accepted by supplier", func(t *testing.T) {

// 	})

// 	t.Run("it 400 if client PUT order to declined by supplier", func(t *testing.T) {

// 	})

// 	t.Run("it 404 if PUT order of non-existing id", func(t *testing.T) {

// 	})

// 	t.Run("it 401 if PUT order unathorized", func(t *testing.T) {

// 	})

// 	t.Run("it 403 if PUT order of not owning user", func(t *testing.T) {

// 	})
// }

// func asserItemNotInModel(t *testing.T, itemModel *data.StubItemModel, itemId int64) {
// 	_, err := itemModel.GetById(itemId)
// 	if !errors.Is(err, data.ErrRecordNotFound) {
// 		t.Fatalf("Wanted to have not found error, got %v", err)
// 	}
// }

// func assertItemResponse(t *testing.T, response *httptest.ResponseRecorder, want data.Item) {
// 	var got data.Item
// 	json.NewDecoder(response.Body).Decode(&got)
// 	if got != want {
// 		t.Fatalf("In response: want %v, got %v", want, got)
// 	}
// }

func createPostOrderRequest(t *testing.T, dto data.PostOrderDto, clientId int64) *http.Request {
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(dto)
	request, err := http.NewRequest(http.MethodPost, "/v1/orders", requestBody)
	tester.AssertNoError(t, err)
	request.Header.Set("Authorization", "Bearer "+strings.Repeat(strconv.FormatInt(clientId, 10), 26))
	return request
}

func parseOrderResponse(t *testing.T, response *httptest.ResponseRecorder) data.Order {
	var got data.Order
	err := json.NewDecoder(response.Body).Decode(&got)
	tester.AssertNoError(t, err)
	return got
}

func assertOrder(t *testing.T, got, want data.Order) {
	if got.Id != want.Id {
		t.Fatalf("Expected order with id %v, got %v", want.Id, got.Id)
	}
	if got.StateId != want.StateId {
		t.Fatalf("Expected order with state id %v, got %v", want.StateId, got.StateId)
	}
	if !data.EqualArrays(got.ItemIds, want.ItemIds) {
		t.Fatalf("Expected order with item ids %v, got %v", want.ItemIds, got.ItemIds)
	}
	if got.MessageId != want.MessageId {
		t.Fatalf("Expected order with message id %v, got %v", want.MessageId, got.MessageId)
	}
}

func assertOrderInModel(t *testing.T, orderModel *data.StubOrderModel, orderId int64, want data.Order) {
	got, err := orderModel.GetById(orderId)
	tester.AssertNoError(t, err)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("In Item Model: want %v, got %v", want, got)
	}
}
