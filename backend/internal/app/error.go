package app

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrUnathorized = errors.New("Unathorized")

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func (a *Application) logError(r *http.Request, err error) {
	a.logger.Print(err)
}

// generic function used by error responses
func (a *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, response ErrorResponse) {
	err := writeJsonResponse(w, status, response, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

func (a *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	a.errorResponse(w, r, http.StatusInternalServerError, ErrorResponse{Message: message})
}

func (a *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	a.errorResponse(w, r, http.StatusNotFound, ErrorResponse{Message: message})
}

func (a *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request, allowedMethods ...string) {
	w.Header().Set("Allow", strings.Join(allowedMethods, "; "))
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	a.errorResponse(w, r, http.StatusMethodNotAllowed, ErrorResponse{Message: message})
}

func (a *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.errorResponse(w, r, http.StatusBadRequest, ErrorResponse{Message: err.Error()})
}

func (a *Application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	a.errorResponse(w, r, http.StatusUnprocessableEntity, ErrorResponse{Message: "Validation error", Errors: errors})
}

func (a *Application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	a.errorResponse(w, r, http.StatusConflict, ErrorResponse{Message: message})
}

func (a *Application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	a.errorResponse(w, r, http.StatusUnauthorized, ErrorResponse{Message: message})
}

func (a *Application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	a.errorResponse(w, r, http.StatusUnauthorized, ErrorResponse{Message: message})
}
