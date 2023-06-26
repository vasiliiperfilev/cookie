package app

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
)

func (a *Application) handlePostItem(w http.ResponseWriter, r *http.Request) {
	user, err := a.AuthenticateHttpRequest(w, r)
	if err != nil {
		a.invalidCredentialsResponse(w, r)
		return
	}
	if user.Type != data.SupplierUserType {
		a.forbiddenResponse(w, r)
		return
	}
	var dto data.PostItemDto
	err = readJsonFromBody(w, r, &dto)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidatePostItemInput(v, dto); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	item := data.Item{
		SupplierId: user.Id,
		Unit:       dto.Unit,
		Size:       dto.Size,
		Name:       dto.Name,
		ImageUrl:   dto.ImageUrl,
	}
	err = a.models.Item.Insert(&item)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	writeJsonResponse(w, http.StatusCreated, item, nil)
}

func (a *Application) handleGetItem(w http.ResponseWriter, r *http.Request) {
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
	itemId, _ := strconv.ParseInt(getField(r, 0), 10, 64)
	item, err := a.models.Item.GetById(itemId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	writeJsonResponse(w, http.StatusOK, item, nil)
}

func (a *Application) handleGetAllItems(w http.ResponseWriter, r *http.Request) {
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
	supplierId, err := strconv.ParseInt(r.URL.Query().Get("supplierId"), 10, 64)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	items, err := a.models.Item.GetAllBySupplierId(supplierId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	writeJsonResponse(w, http.StatusOK, items, nil)
}
