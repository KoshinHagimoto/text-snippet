package main

import (
	"log"
	"net/http"
	"text-snippet/app/dao"
	handler "text-snippet/app/handlers"
	"text-snippet/app/middleware"

	"github.com/gorilla/mux"
)

var snipDao *dao.SnippetDAO

func main() {
	var err error
	snipDao, err = dao.InitSnippetDAO()
	if err != nil {
		log.Fatal(err)
	}
	defer snipDao.Finalize()

	r := mux.NewRouter()

	r.HandleFunc("/snippet", middleware.CORSMiddleware(handler.CreateSnippetHandler(snipDao))).Methods("POST")
	r.HandleFunc("/snippets", middleware.CORSMiddleware(handler.GetSnippetListHandler(snipDao))).Methods("GET")

	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/").Handler(fs)

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
