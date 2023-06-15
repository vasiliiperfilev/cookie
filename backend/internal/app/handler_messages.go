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
		a.handleGetMessages(w, r)
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

// handles /v1/conversations/([0-9]+)/messages route
func (a *Application) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	conversationId, _ := strconv.ParseInt(getField(r, 0), 10, 64)

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
