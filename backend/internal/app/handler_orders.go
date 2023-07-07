package app

import (
	"errors"
	"net/http"
	"strconv"

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
	dto.ClientId = user.Id
	v := validator.New()
	if data.ValidatePostOrderInput(v, dto); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	order, err := a.repositories.Order.Insert(dto)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUnprocessableEntity):
			v.AddError("itemIds", "At least one of order items doesn't exist")
			a.failedValidationResponse(w, r, v.Errors)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	writeJsonResponse(w, http.StatusCreated, order, nil)
}

func (a *Application) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	_, err := a.AuthenticateHttpRequest(w, r)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnathorized):
			a.invalidAuthenticationTokenResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	orderId, _ := strconv.ParseInt(getField(r, 0), 10, 64)
	order, err := a.models.Order.GetById(orderId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	writeJsonResponse(w, http.StatusOK, order, nil)
}
