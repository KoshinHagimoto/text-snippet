package handler

import (
	"encoding/json"
	"net/http"
	"text-snippet/app/dao"
)

func GetSnippetListHandler(snipDao *dao.SnippetDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snippets, err := snipDao.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(snippets)
	}
}
