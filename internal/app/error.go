package app

import (
	"fmt"
	"net/http"
	"strings"
)

func (a *Application) logError(r *http.Request, err error) {
	a.logger.Print(err)
}

// generic function used by error responses
func (a *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	data := map[string]any{"error": message}

	err := writeJSON(w, status, data, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

func (a *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	a.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (a *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	a.errorResponse(w, r, http.StatusNotFound, message)
}

func (a *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request, allowedMethods ...string) {
	w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	a.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (a *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	a.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
