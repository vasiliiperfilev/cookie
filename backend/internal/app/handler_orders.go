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
	dto.ClientId = user.Id
	v := validator.New()
	if data.ValidatePostOrderInput(v, dto); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	order, err := a.repositories.Order.Insert(dto)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
	writeJsonResponse(w, http.StatusCreated, order, nil)
}
