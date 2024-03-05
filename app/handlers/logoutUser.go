package handler

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func LogoutUserHandler(sessionStore *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessionStore.Get(r, "login-session")

		if _, ok := session.Values["user_id"]; !ok {
			http.Error(w, "User not logged in", http.StatusBadRequest)
			return
		}

		// セッションデータを削除
		session.Options.MaxAge = -1

		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
