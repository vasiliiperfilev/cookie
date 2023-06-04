package app

import "net/http"

func (a *Application) conversationsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, nil, nil)
}
