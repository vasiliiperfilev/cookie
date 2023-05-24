package app

import (
	"net/http"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
)

func (a *Application) authRegisterHandler(w http.ResponseWriter, r *http.Request) {
	registerUserInput := new(data.RegisterUserInput)
	err := readJSON(w, r, registerUserInput)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateRegisterUserInput(v, registerUserInput); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := data.User{Email: registerUserInput.Email, Type: registerUserInput.Type, ImageId: registerUserInput.ImageId}
	writeJSON(w, http.StatusOK, user, nil)
}
