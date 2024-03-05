package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/sessions"
)

func GetUserStatusHandler(sessionStore *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := sessionStore.Get(r, "login-session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		loggedIn := session.Values["user_id"] != nil
		json.NewEncoder(w).Encode(map[string]bool{"loggedIn": loggedIn})
	}
}
