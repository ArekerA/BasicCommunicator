package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.next.ServeHTTP(w, r)
}
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := ""
	provider := ""
	if len(segs) >= 4 {
		action = segs[2]
		provider = segs[3]
	}
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, fmt.Sprintf("Błąd podczas próby pobrania dostawcy %s %s", provider, err), http.StatusBadRequest)
			return
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Błąd wywołania GetBeginAuthURL %s %s", provider, err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Akcja autoryzacyjna %s nie jest obsługiwana", action)
	}

}
