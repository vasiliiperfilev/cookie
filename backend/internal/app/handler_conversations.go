package app

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
	"golang.org/x/exp/slices"
)

func (a *Application) handlePostConversation(w http.ResponseWriter, r *http.Request) {
	user, err := a.AuthenticateHttpRequest(w, r)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnathorized):
			a.invalidAuthenticationTokenResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	dto := data.PostConversationDto{}
	err = readJsonFromBody(w, r, &dto)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	if !slices.Contains(dto.UserIds, user.Id) {
		a.badRequestResponse(w, r, ErrUnathorized)
		return
	}
	v := validator.New()
	c, err := a.models.Conversation.Insert(dto)
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

	writeJsonResponse(w, http.StatusCreated, c, nil)
}

func (a *Application) handleGetConversation(w http.ResponseWriter, r *http.Request) {
	user, err := a.AuthenticateHttpRequest(w, r)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnathorized):
			a.invalidAuthenticationTokenResponse(w, r)
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

	writeJsonResponse(w, http.StatusOK, conversations, nil)
}
