package app

import (
	"net/http"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
)

func (a *Application) handlePostOrder(w http.ResponseWriter, r *http.Request) {
	user, err := a.AuthenticateHttpRequest(w, r)
	if err != nil {
		a.invalidCredentialsResponse(w, r)
		return
	}
	var dto data.PostOrderDto
	err = readJsonFromBody(w, r, &dto)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidatePostOrderInput(v, dto); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	message := data.Message{
		ConversationId: dto.ConversationId,
		Content:        "Order created",
		SenderId:       user.Id,
	}
	err = a.models.Message.Insert(&message)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	order := data.Order{
		ItemIds:   dto.ItemIds,
		StateId:   data.OrderStateCreated,
		MessageId: message.Id,
	}
	order, err = a.models.Order.Insert(order)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	writeJsonResponse(w, http.StatusCreated, order, nil)
}
