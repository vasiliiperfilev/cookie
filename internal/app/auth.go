package app

import (
	"net/http"

	"github.com/vasiliiperfilev/cookie/internal/data"
)

func (a *Application) authRegisterHandler(w http.ResponseWriter, r *http.Request) {
	registerUserInput := new(data.RegisterUserInput)

	err := readJSON(w, r, registerUserInput)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	registerUserResponse := data.RegisterUserResponse{Email: registerUserInput.Email, Type: registerUserInput.Type, ImageId: registerUserInput.ImageId}
	writeJSON(w, http.StatusOK, registerUserResponse, nil)
}
