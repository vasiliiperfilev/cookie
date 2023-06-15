package app

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/vasiliiperfilev/cookie/internal/data"
)

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
