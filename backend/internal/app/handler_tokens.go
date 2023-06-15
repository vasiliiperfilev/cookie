package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
)

type HandlerTokenResponse struct {
	User  *data.User  `json:"user"`
	Token *data.Token `json:"token"`
}

func (a *Application) tokensHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.handlePostToken(w, r)
	case http.MethodOptions:
		w.Header().Set("Allow", http.MethodPost)
		err := writeJsonResponse(w, http.StatusOK, nil, nil)
		if err != nil {
			a.serverErrorResponse(w, r, err)
		}
	default:
		a.methodNotAllowedResponse(w, r, http.MethodPost)
	}
}

func (a *Application) handlePostToken(w http.ResponseWriter, r *http.Request) {
	// Parse the email and password from the request body.
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := readJsonFromBody(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Validate the email provided by the client.
	v := validator.New()
	data.ValidateEmail(v, input.Email)
	if !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := a.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.invalidCredentialsResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Check if the provided password matches the actual password for the user.
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// If the passwords don't match, then we call the app.invalidCredentialsResponse()
	// helper again and return.
	if !match {
		a.invalidCredentialsResponse(w, r)
		return
	}

	// Otherwise, if the password is correct, we generate a new token with a 24-hour
	// expiry time and the scope 'authentication'.
	token, err := a.models.Token.New(user.Id, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	err = writeJsonResponse(w, http.StatusCreated, HandlerTokenResponse{User: user, Token: token}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
