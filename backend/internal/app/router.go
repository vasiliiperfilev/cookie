package app

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

func (a *Application) routes() http.Handler {
	routes := []route{
		newRoute(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler),
		newRoute(http.MethodPost, "/v1/users", a.handlePostUser),
		newRoute(http.MethodGet, "/v1/users", a.handleGetUsers),
		newRoute(http.MethodPost, "/v1/tokens", a.handlePostToken),
		newRoute(http.MethodPost, "/v1/conversations", a.handlePostConversation),
		newRoute(http.MethodGet, "/v1/conversations", a.handleGetConversation),
		newRoute(http.MethodGet, "/v1/chat", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			a.wsChatHandler(a.hub, w, r)
		})),
		newRoute(http.MethodGet, "/v1/conversations/([0-9]+)/messages", a.handleGetMessages),
		newRoute(http.MethodGet, "/v1/messages/([0-9]+)", a.handleGetMessage),
		newRoute(http.MethodPost, "/v1/items", a.handlePostItem),
		newRoute(http.MethodGet, "/v1/items", a.handleGetAllItems),
		newRoute(http.MethodGet, "/v1/items/([0-9]+)", a.handleGetItem),
		newRoute(http.MethodPut, "/v1/items/([0-9]+)", a.handlePutItem),
		newRoute(http.MethodDelete, "/v1/items/([0-9]+)", a.handleDeleteItem),
		newRoute(http.MethodPost, "/v1/orders", a.handlePostOrder),
		newRoute(http.MethodGet, "/v1/orders", a.handleGetAllOrders),
		newRoute(http.MethodGet, "/v1/orders/([0-9]+)", a.handleGetOrder),
		newRoute(http.MethodPatch, "/v1/orders/([0-9]+)", a.handlePatchOrder),
	}
	return NewRouter(routes)
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func NewRouter(routes []route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allow []string
		for _, route := range routes {
			matches := route.regex.FindStringSubmatch(r.URL.Path)
			if len(matches) > 0 {
				if r.Method != route.method {
					allow = append(allow, route.method)
					continue
				}
				ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
				route.handler(w, r.WithContext(ctx))
				return
			}
		}
		if len(allow) > 0 && r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allow, ", "))
			writeJsonResponse(w, http.StatusOK, nil, nil)
			return
		}
		if len(allow) > 0 {
			w.Header().Set("Allow", strings.Join(allow, ", "))
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.NotFound(w, r)
	})
}

type ctxKey struct{}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}
