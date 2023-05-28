package app

import (
	"errors"
	"net/http"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
)

func (a *Application) userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostUser(w, r, a)
	default:
		a.methodNotAllowedResponse(w, r, http.MethodPost)
	}
}

func handlePostUser(w http.ResponseWriter, r *http.Request, a *Application) {
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
	err = user.Password.Set(registerUserInput.Password)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = a.models.User.Insert(&user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			a.failedValidationResponse(w, r, v.Errors)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	writeJSON(w, http.StatusOK, user, nil)
}
