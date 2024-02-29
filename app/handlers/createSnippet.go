package handler

import (
	"encoding/json"
	"net/http"
	"text-snippet/app/dao"
	"text-snippet/app/object"
)

func CreateSnippetHandler(snipDao *dao.SnippetDAO) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var snippet object.Snippet
		if err := json.NewDecoder(r.Body).Decode(&snippet); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := snipDao.Save(snippet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]int{"id": id})
	}
}
