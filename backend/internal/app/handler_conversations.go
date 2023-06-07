package app

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
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
		err := writeJsonResponse(w, http.StatusOK, nil, nil)
		if err != nil {
			a.serverErrorResponse(w, r, err)
		}
	default:
		a.methodNotAllowedResponse(w, r, http.MethodPost, http.MethodGet)
	}
}

func handlePostConversation(w http.ResponseWriter, r *http.Request, a *Application) {
	conversation := new(data.Conversation)
	err := readJsonFromBody(w, r, conversation)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	err = a.models.Conversation.Insert(*conversation)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateConversation):
			v.AddError("userIds", "conversation with these users already exists")
			a.failedValidationResponse(w, r, v.Errors)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	writeJsonResponse(w, http.StatusOK, conversation, nil)
}

func handleGetConversation(w http.ResponseWriter, r *http.Request, a *Application) {
	userId, err := strconv.ParseInt(r.URL.Query().Get("userId"), 10, 64)
	if err != nil || userId < 1 {
		a.notFoundResponse(w, r)
		return
	}
	conversations, err := a.models.Conversation.GetAllByUserId(userId)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	writeJsonResponse(w, http.StatusOK, conversations, nil)
}
