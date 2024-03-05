package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text-snippet/app/config"
	"text-snippet/app/dao"
	handler "text-snippet/app/handlers"
	"text-snippet/app/middleware"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func main() {
	var err error
	db, err := config.InitDAO()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	CSRFSecret := os.Getenv("CSRF_SECRET")
	if CSRFSecret == "" {
		log.Fatal("CSRF_SECRET is not set")
	}

	// HTTPSを使用していない場合はcsrf.Secure(false)を設定
	CSRFMiddleware := csrf.Protect([]byte(CSRFSecret), csrf.Secure(false))

	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   1800,
		HttpOnly: true,
		Secure:   false, // HTTPSを使用していない場合はfalse
	}

	r := mux.NewRouter()

	snipDao := dao.NewSnippetDAO(db)
	userDao := dao.NewUserDAO(db)

	// CSRFトークンを取得するエンドポイント
	r.HandleFunc("/csrf-token", func(w http.ResponseWriter, r *http.Request) {
		csrfToken := csrf.Token(r)                                           // CSRFトークンを取得
		json.NewEncoder(w).Encode(map[string]string{"csrfToken": csrfToken}) // トークンをJSONで返す
	}).Methods("GET")

	r.HandleFunc("/snippet", middleware.CORSMiddleware(handler.CreateSnippetHandler(snipDao))).Methods("POST")
	r.HandleFunc("/snippets", middleware.CORSMiddleware(handler.GetSnippetListHandler(snipDao))).Methods("GET")
	r.HandleFunc(("/user/register"), middleware.CORSMiddleware(handler.RegisterUserHandler(userDao))).Methods("POST")
	r.HandleFunc("/user/login", middleware.CORSMiddleware(handler.LoginUserHandler(userDao, sessionStore))).Methods("POST")
	r.HandleFunc("/user/logout", middleware.CORSMiddleware(handler.LogoutUserHandler(sessionStore))).Methods("POST")

	http.Handle("/", CSRFMiddleware(r))

	fs := http.FileServer(http.Dir("./public"))
	r.PathPrefix("/").Handler(fs)

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", CSRFMiddleware(r))
}
