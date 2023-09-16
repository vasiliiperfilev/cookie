package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	userModel := data.NewStubUserModel(generateUsers(4))
	conversationModel := data.NewStubConversationModel(generateConversation(4), userModel)
	messageModel := data.NewStubMessageModel(generateConversation(4), []data.Message{{Id: 1, ConversationId: 1, PrevMessageId: 0}})
	orderModel := data.NewStubOrderModel([]data.Order{}, itemModel, conversationModel, messageModel)
	models := data.Models{
		Conversation: conversationModel,
		User:         userModel,
		Item:         itemModel,
		Message:      messageModel,
		Order:        orderModel,
	}
	server := app.New(cfg, logger, models)

	t.Run("it POST order with correct values", func(t *testing.T) {
		clientId := int64(1)
		dto := data.PostOrderDto{
			ConversationId: 1,
			Items: []data.ItemQuantity{
				{
					ItemId:   1,
					Quantity: 1,
				},
				{
					ItemId:   2,
					Quantity: 1,
				},
			},
		}
		request := createPostOrderRequest(t, dto, clientId)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		messages, err := messageModel.GetAllByConversationId(1)
		tester.AssertNoError(t, err)
		want := data.Order{
			Items: []data.ItemQuantity{
				{
					ItemId:   1,
					Quantity: 1,
				},
				{
					ItemId:   2,
					Quantity: 1,
				},
			},
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
			Items: []data.ItemQuantity{
				{
					ItemId:   4,
					Quantity: 1,
				},
				{
					ItemId:   5,
					Quantity: 1,
				},
			},
		}
		wantMsgCount := countUserMessages(t, messageModel, clientId)
		wantOrderCount := countUserOrder(t, orderModel, clientId)
		request := createPostOrderRequest(t, dto, clientId)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		gotOrderCount := countUserOrder(t, orderModel, clientId)
		gotMsgCount := countUserMessages(t, messageModel, clientId)
		tester.AssertValue(t, gotMsgCount, wantMsgCount, "Expected to not have new messages")
		tester.AssertValue(t, gotOrderCount, wantOrderCount, "Expected to not have new orders")
		tester.AssertStatus(t, response.Code, http.StatusUnprocessableEntity)
		assertContentType(t, response, app.JsonContentType)
	})

	t.Run("it 422 if POST order with items of different suppliers", func(t *testing.T) {
		// TODO: implement
	})

	t.Run("it 401 if POST order unathorized", func(t *testing.T) {
		// TODO: implement
	})

	t.Run("it 403 if POST order to not own conversation", func(t *testing.T) {
		// TODO: implement
	})
}

func TestOrderGet(t *testing.T) {
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
	userModel := data.NewStubUserModel(generateUsers(4))
	conversationModel := data.NewStubConversationModel(generateConversation(4), userModel)
	messageModel := data.NewStubMessageModel(generateConversation(4), []data.Message{{Id: 1, ConversationId: 1, PrevMessageId: 0}})
	orderModel := data.NewStubOrderModel([]data.Order{}, itemModel, conversationModel, messageModel)
	models := data.Models{
		Conversation: conversationModel,
		User:         userModel,
		Item:         itemModel,
		Message:      messageModel,
		Order:        orderModel,
	}
	server := app.New(cfg, logger, models)

	t.Run("it GET order", func(t *testing.T) {
		dto := data.PostOrderDto{
			ClientId:       1,
			ConversationId: 1,
			Items: []data.ItemQuantity{
				{
					ItemId:   1,
					Quantity: 1,
				},
				{
					ItemId:   2,
					Quantity: 1,
				},
			},
		}
		want, err := orderModel.Insert(dto)
		tester.AssertNoError(t, err)

		request := createGetOrderRequest(t, 1, 1)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		tester.AssertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		got := parseOrderResponse(t, response)
		assertOrder(t, got, want)
	})

	t.Run("it 404 if GET non-existing order", func(t *testing.T) {
		request := createGetOrderRequest(t, 123, 1)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		tester.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("it 401 if GET order unathorized", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/orders/%v", 1), nil)
		tester.AssertNoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		tester.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("it 403 if GET not own order", func(t *testing.T) {
		// TODO: implement
	})
}

func TestOrderGetAll(t *testing.T) {
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
	userModel := data.NewStubUserModel(generateUsers(4))
	conversationModel := data.NewStubConversationModel(generateConversation(4), userModel)
	messageModel := data.NewStubMessageModel(generateConversation(4), []data.Message{{Id: 1, ConversationId: 1, PrevMessageId: 0}})
	orderModel := data.NewStubOrderModel([]data.Order{}, itemModel, conversationModel, messageModel)
	models := data.Models{
		Conversation: conversationModel,
		User:         userModel,
		Item:         itemModel,
		Message:      messageModel,
		Order:        orderModel,
	}
	server := app.New(cfg, logger, models)

	t.Run("it GET all orders of own id", func(t *testing.T) {
		dto := data.PostOrderDto{
			ClientId:       1,
			ConversationId: 1,
			Items: []data.ItemQuantity{
				{
					ItemId:   1,
					Quantity: 1,
				},
				{
					ItemId:   2,
					Quantity: 1,
				},
			},
		}
		client, err := userModel.GetById(1)
		tester.AssertNoError(t, err)
		order1, err := orderModel.Insert(dto)
		tester.AssertNoError(t, err)
		order2, err := orderModel.Insert(dto)
		tester.AssertNoError(t, err)
		order1.Client = client
		order2.Client = client
		want := []data.Order{order1, order2}

		request := createGetAllOrdersRequest(t, 1)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		tester.AssertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		got := tester.ParseResponse[[]data.Order](t, response)
		for _, order := range want {
			assertOrderInArray(t, order, got)
		}
	})

	t.Run("it 404 if GET all orders of non-existing", func(t *testing.T) {

	})

	t.Run("it 401 if GET all orders unathorized", func(t *testing.T) {

	})

	t.Run("it 403 if GET all orders of not owning user", func(t *testing.T) {

	})
}

func TestOrderPatch(t *testing.T) {
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
	testOrder := data.Order{
		Id: 1,
		Items: []data.ItemQuantity{
			{
				ItemId:   1,
				Quantity: 1,
			},
			{
				ItemId:   2,
				Quantity: 1,
			},
		},
		StateId:   data.OrderStateCreated,
		MessageId: 1,
	}
	userModel := data.NewStubUserModel(generateUsers(4))
	conversationModel := data.NewStubConversationModel(generateConversation(4), userModel)
	messageModel := data.NewStubMessageModel(generateConversation(4), []data.Message{{Id: 1, ConversationId: 1, PrevMessageId: 0, SenderId: 1}})
	orderModel := data.NewStubOrderModel([]data.Order{testOrder}, itemModel, conversationModel, messageModel)
	models := data.Models{
		Conversation: conversationModel,
		User:         userModel,
		Item:         itemModel,
		Message:      messageModel,
		Order:        orderModel,
	}
	server := app.New(cfg, logger, models)

	// Order states with owner:
	// Created - 1 (client)
	// Accepted - 2 (supplier)
	// Declined - 3 (supplier)
	// Fulfilled - 4 (supplier)
	// Confirmed fulfillment - 5 (client)
	// Supplier changes - 6 (supplier)
	// Client changes - 7 (client)
	// switch handling function depending on what was changed: items or state
	// switch patch state and require permission depending on state
	// validate patch order inputs

	t.Run("it 200 if supplier PATCH order state to accepted", func(t *testing.T) {
		supplierId := int64(2)
		dto := data.PatchOrderDto{
			StateId: data.OrderStateAccepted,
		}
		request := createPatchOrderRequest(t, dto, supplierId, 1)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		want := testOrder
		want.StateId = data.OrderStateAccepted
		tester.AssertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, app.JsonContentType)
		got := tester.ParseResponse[data.Order](t, response)
		assertOrder(t, got, want)
		assertOrderInModel(t, orderModel, got.Id, want)
	})

	// 	t.Run("it 200 if supplier PATCH order state to fulfilled", func(t *testing.T) {
	// 		// stop at this point and implement frontend
	// 	})

	// 	t.Run("it 200 if client PATCH order state to confirm fulfielment", func(t *testing.T) {

	// 	})

	// 	t.Run("it 200 if supplier PATCH order state to suggest changes before fulfielment", func(t *testing.T) {

	// 	})

	// 	t.Run("it 200 if client PATCH order state to suggest changes before fulfielment", func(t *testing.T) {

	// 	})

	// 	t.Run("it 200 if client PATCH order state to accept supplier changes", func(t *testing.T) {

	// 	})

	// 	t.Run("it 200 if client PATCH order state to accept client changes", func(t *testing.T) {

	// 	})

	// 	t.Run("it 200 if supplier PATCH order state to decline before accepting/suggesting changes", func(t *testing.T) {

	// 	})

	// 	t.Run("it 400 if supplier PATCH order state to decline after fulfilled", func(t *testing.T) {

	// 	})

	// 	t.Run("it 400 if supplier PATCH order state to accepted after own changes", func(t *testing.T) {

	// 	})

	// 	t.Run("it 400 if client PATCH order state to accepted after own changes", func(t *testing.T) {

	// 	})

	// 	t.Run("it 400 if supplier PATCH order state to canceled by client", func(t *testing.T) {

	// 	})

	// 	t.Run("it 400 if client PATCH order state to accepted by supplier", func(t *testing.T) {

	// 	})

	// 	t.Run("it 400 if client PATCH order state to declined by supplier", func(t *testing.T) {

	// 	})

	// 	t.Run("it 404 if PATCH order state of non-existing id", func(t *testing.T) {

	// 	})

	// 	t.Run("it 401 if PATCH order state unathorized", func(t *testing.T) {

	// 	})

	// 	t.Run("it 403 if PATCH order state of not owning user", func(t *testing.T) {

	// })
}

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

func createPatchOrderRequest(t *testing.T, dto data.PatchOrderDto, clientId int64, orderId int64) *http.Request {
	requestBody := new(bytes.Buffer)
	json.NewEncoder(requestBody).Encode(dto)
	request, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/v1/orders/%v", orderId), requestBody)
	tester.AssertNoError(t, err)
	request.Header.Set("Authorization", "Bearer "+strings.Repeat(strconv.FormatInt(clientId, 10), 26))
	return request
}

func createGetOrderRequest(t *testing.T, orderId int64, clientId int64) *http.Request {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/orders/%v", orderId), nil)
	tester.AssertNoError(t, err)
	request.Header.Set("Authorization", "Bearer "+strings.Repeat(strconv.FormatInt(clientId, 10), 26))
	return request
}

func createGetAllOrdersRequest(t *testing.T, clientId int64) *http.Request {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/orders?userId=%v", clientId), nil)
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
	if !data.EqualArraysContent(got.Items, want.Items) {
		t.Fatalf("Expected order with item ids %v, got %v", want.Items, got.Items)
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

func countUserMessages(t *testing.T, messageModel *data.StubMessageModel, userId int64) int {
	messages, err := messageModel.GetAllByConversationId(1)
	tester.AssertNoError(t, err)
	return len(messages)
}

func countUserOrder(t *testing.T, orderModel *data.StubOrderModel, userId int64) int {
	orders, err := orderModel.GetAllByUserId(userId)
	tester.AssertNoError(t, err)
	return len(orders)
}

func assertOrderInArray(t *testing.T, o data.Order, arr []data.Order) {
	for _, order := range arr {
		if reflect.DeepEqual(order, o) {
			return
		}
	}
	t.Fatalf("Expected to find order %v in array %v", o, arr)
}
