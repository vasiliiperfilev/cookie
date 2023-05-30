package app

import (
	"net/http"
)

func (a *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := writeJSON(w, http.StatusOK, a.GetState(), nil)
		if err != nil {
			a.serverErrorResponse(w, r, err)
		}
	default:
		a.methodNotAllowedResponse(w, r, http.MethodGet)
	}
}
