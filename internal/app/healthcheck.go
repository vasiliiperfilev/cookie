package app

import (
	"net/http"
)

func (a *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := writeJSON(w, http.StatusOK, a.GetState(), nil)
		if err != nil {
			a.logger.Print(err)
			http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		}
	default:
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
	}
}
