package app

import (
	"net/http"

	"github.com/vasiliiperfilev/cookie/internal/data"
)

func (a *Application) handlePostItem(w http.ResponseWriter, r *http.Request) {
	user, err := a.AuthenticateHttpRequest(w, r)
	if err != nil {
		a.invalidCredentialsResponse(w, r)
		return
	}
	var dto data.PostItemDto
	err = readJsonFromBody(w, r, &dto)
	if err != nil {
		a.badRequestResponse(w, r, err)
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
