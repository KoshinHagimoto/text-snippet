package handler

import (
	"encoding/json"
	"net/http"
	"text-snippet/app/dao"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func LoginUserHandler(userDao *dao.UserDAO, sessionStore *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials struct {
			Username string
			Password string
		}

		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user, err := userDao.GetUserByUsername(credentials.Username)
		if err != nil {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
			return
		}

		session, _ := sessionStore.Get(r, "login-session")
		session.Values["user_id"] = user.ID
		session.Save(r, w)

		json.NewEncoder(w).Encode("Logged in successfully")
	}
}
