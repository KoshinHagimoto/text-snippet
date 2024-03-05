package handler

import (
	"encoding/json"
	"net/http"
	"text-snippet/app/config"
	"text-snippet/app/dao"
	"text-snippet/app/object"

	"regexp"

	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

// EmailRegex represents the regular expression pattern for email validation
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func RegisterUserHandler(userDao *dao.UserDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user object.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate email
		if !EmailRegex.MatchString(user.Email) {
			http.Error(w, "Invalid email", http.StatusBadRequest)
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

		// Email verification
		verificationURL := config.GenerateSignedURL(user.ID, user.Email)
		config.SendVerificationEmail(user.Email, verificationURL)

		// Set the X-CSRF-Token header
		csrfToken := csrf.Token(r)
		w.Header().Set("X-CSRF-Token", csrfToken)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "user created"})
	}
}
