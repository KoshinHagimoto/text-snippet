package main

import (
	"log"
	"net/http"
	"text-snippet/app/config"
	"text-snippet/app/dao"
	handler "text-snippet/app/handlers"
	"text-snippet/app/middleware"

	"github.com/gorilla/mux"
)

func main() {
	var err error
	db, err := config.InitDAO()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()

	snipDao := dao.NewSnippetDAO(db)
	userDao := dao.NewUserDAO(db)

	r.HandleFunc("/snippet", middleware.CORSMiddleware(handler.CreateSnippetHandler(snipDao))).Methods("POST")
	r.HandleFunc("/snippets", middleware.CORSMiddleware(handler.GetSnippetListHandler(snipDao))).Methods("GET")
	r.HandleFunc(("/user/register"), middleware.CORSMiddleware(handler.RegisterUserHandler(userDao))).Methods("POST")

	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/").Handler(fs)

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
