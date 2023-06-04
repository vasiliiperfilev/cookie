package app

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vasiliiperfilev/cookie/internal/data"
)

func (a *Application) conversationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostConversation(w, r, a)
	case http.MethodGet:
		handleGetConversation(w, r, a)
	case http.MethodOptions:
		allowed := []string{http.MethodPost, http.MethodGet}
		w.Header().Set("Allow", strings.Join(allowed, "; "))
		err := writeJSON(w, http.StatusOK, nil, nil)
		if err != nil {
			a.serverErrorResponse(w, r, err)
		}
	default:
		a.methodNotAllowedResponse(w, r, http.MethodPost, http.MethodGet)
	}
}

func handlePostConversation(w http.ResponseWriter, r *http.Request, a *Application) {
	conversation := new(data.Conversation)
	err := readJSON(w, r, conversation)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	conversation.Id = 1
	conversation.LastMessageId = -1

	writeJSON(w, http.StatusOK, conversation, nil)
}

func handleGetConversation(w http.ResponseWriter, r *http.Request, a *Application) {
	id, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil || id < 1 {
		a.notFoundResponse(w, r)
		return
	}
	conversations := []data.Conversation{
		{
			Id:            1,
			UserIds:       []int64{1, 2},
			LastMessageId: -1,
		},
	}
	writeJSON(w, http.StatusOK, conversations, nil)
}
