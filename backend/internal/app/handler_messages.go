package app

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/vasiliiperfilev/cookie/internal/data"
)

func (a *Application) messagesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetMessages(w, r, a)
	case http.MethodOptions:
		w.Header().Set("Allow", http.MethodGet)
		err := writeJsonResponse(w, http.StatusOK, nil, nil)
		if err != nil {
			a.serverErrorResponse(w, r, err)
		}
	default:
		a.methodNotAllowedResponse(w, r, http.MethodPost, http.MethodGet)
	}
}

func handleGetMessages(w http.ResponseWriter, r *http.Request, a *Application) {
	conversationId, err := extractConversationId(r.URL.Path)
	if err != nil || conversationId < 1 {
		a.notFoundResponse(w, r)
		return
	}

	messages, err := a.models.Message.GetAllByConversationId(conversationId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	writeJsonResponse(w, http.StatusOK, messages, nil)
}

func extractConversationId(uri string) (int64, error) {
	noPreffix := strings.TrimPrefix(uri, "/v1/conversations/")
	idStr := strings.TrimSuffix(noPreffix, "/messages")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
