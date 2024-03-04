package handler

import (
	"encoding/json"
	"net/http"
	"text-snippet/app/dao"
	"text-snippet/app/object"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUserHandler(userDao *dao.UserDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user object.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user.PasswordHash = string(hashedPassword)
		if err := userDao.CreateUser(&user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "user created"})
	}
}
