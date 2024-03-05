package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"text-snippet/app/dao"
	"text-snippet/app/middleware"
)

func VerifyEmailHandler(userDao *dao.UserDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if middleware.VerifySignature(r) {
			// URL署名が正しい場合の処理
			userId, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}
			err = userDao.UpdateEmailVerified(userId, true)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintln(w, "Email verified successfully")
		} else {
			http.Error(w, "Invalid signature", http.StatusBadRequest)
			return
		}
	}
}
