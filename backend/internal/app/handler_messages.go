package app

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"golang.org/x/exp/slices"
)

// handles /v1/conversations/([0-9]+)/messages route
func (a *Application) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	user, err := a.AuthenticateHttpRequest(w, r)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnathorized):
			a.invalidAuthenticationTokenResponse(w, r)
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	conversationId, _ := strconv.ParseInt(getField(r, 0), 10, 64)
	cvs, err := a.models.Conversation.GetById(conversationId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	if !slices.Contains(cvs.UserIds, user.Id) {
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

func (a *Application) handleGetMessage(w http.ResponseWriter, r *http.Request) {
	_, err := a.AuthenticateHttpRequest(w, r)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnathorized):
			a.invalidAuthenticationTokenResponse(w, r)
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	messageId, _ := strconv.ParseInt(getField(r, 0), 10, 64)
	msg, err := a.models.Message.GetById(messageId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	writeJsonResponse(w, http.StatusOK, msg, nil)
}
