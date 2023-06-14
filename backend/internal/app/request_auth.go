package app

import (
	"net/http"
	"strings"

	"github.com/vasiliiperfilev/cookie/internal/data"
	"github.com/vasiliiperfilev/cookie/internal/validator"
)

func (a *Application) AuthenticateHttpRequest(w http.ResponseWriter, r *http.Request) (*data.User, error) {
	// Add the "Vary: Authorization" header to the response. This indicates to any
	// caches that the response may vary based on the value of the Authorization
	// header in the request.
	w.Header().Add("Vary", "Authorization")

	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		return nil, ErrUnathorized
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, ErrUnathorized
	}

	token := headerParts[1]

	v := validator.New()
	if data.ValidateTokenPlaintext(v, token); !v.Valid() {
		return nil, ErrUnathorized
	}

	user, err := a.models.User.GetForToken(data.ScopeAuthentication, token)
	if err != nil {
		return nil, data.ErrRecordNotFound
	}

	return user, nil
}

func (a *Application) AuthenticateWsUpgradeRequest(w http.ResponseWriter, r *http.Request) (*data.User, error) {
	token := r.URL.Query().Get("token")

	v := validator.New()
	if data.ValidateTokenPlaintext(v, token); !v.Valid() {
		return nil, ErrUnathorized
	}

	user, err := a.models.User.GetForToken(data.ScopeAuthentication, token)
	if err != nil {
		return nil, data.ErrRecordNotFound
	}

	return user, nil
}
