package app

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
	"golang.org/x/exp/slices"
)

type ExpandedMessage struct {
	data.Message
	Sender data.User `json:"sender"`
}

type ExpandedConversation struct {
	data.Conversation
	LastMessage ExpandedMessage `json:"lastMessage"`
}

func (a *Application) handlePostConversation(w http.ResponseWriter, r *http.Request) {
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
	conversation := new(data.Conversation)
	err = readJsonFromBody(w, r, conversation)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	if !slices.Contains(conversation.UserIds, user.Id) {
		a.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	err = a.models.Conversation.Insert(conversation)
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

	writeJsonResponse(w, http.StatusCreated, conversation, nil)
}

func (a *Application) handleGetConversation(w http.ResponseWriter, r *http.Request) {
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
	userId, err := strconv.ParseInt(r.URL.Query().Get("userId"), 10, 64)
	if err != nil || userId < 1 || userId != user.Id {
		a.notFoundResponse(w, r)
		return
	}
	expanded := r.URL.Query().Get("expanded")

	conversations, err := a.models.Conversation.GetAllByUserId(userId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	if expanded == "" {
		writeJsonResponse(w, http.StatusOK, conversations, nil)
	} else {
		expandedConvs := []ExpandedConversation{}
		for _, c := range conversations {
			msg, err := a.models.Message.GetById(c.LastMessageId)
			if err != nil {
				a.serverErrorResponse(w, r, err)
				return
			}
			usr, err := a.models.User.GetById(msg.SenderId)
			if err != nil {
				a.serverErrorResponse(w, r, err)
				return
			}
			expandedConvs = append(expandedConvs, ExpandedConversation{LastMessage: ExpandedMessage{Message: *msg, Sender: *usr}, Conversation: c})
		}
		writeJsonResponse(w, http.StatusOK, expandedConvs, nil)
	}
}
