package app

import (
	"encoding/json"
	"net/http"
)

func (a *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", JsonContentType)
	json.NewEncoder(w).Encode(a.GetState())
	w.WriteHeader(http.StatusOK)
}
