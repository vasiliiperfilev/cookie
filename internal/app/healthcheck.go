package app

import (
	"encoding/json"
	"net/http"
)

func (a *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("content-type", JsonContentType)
		json.NewEncoder(w).Encode(a.GetState())
		w.WriteHeader(http.StatusOK)
	default:
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
	}
}
