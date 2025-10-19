package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"pkg/auth"
	"pkg/database"
	"pkg/handler"
	"pkg/repository"
	"pkg/service"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	templatesPath := os.Getenv("TEMPLATES_PATH")

	tmpl, err := template.ParseGlob(fmt.Sprintf("%s/*.html", templatesPath))
	if err != nil {
		log.Fatal(err)
	}

	if err := os.Remove(dbPath); err != nil {
		log.Fatal(err)
	}

	db, err := database.NewSqliteDatabase(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	ur := repository.NewSqliteUserRepository(db)
	us := service.NewUserService(ur, tmpl)

	tr := repository.NewSqliteTokenRepository(db)
	ts := service.NewTokenService(tr)

	rr := repository.NewSqliteResumeRepository(db)
	rs := service.NewResumeService(rr)

	ar := auth.NewAuthRepository()

	router := mux.NewRouter()
	profileHandler := handler.NewProfileHandler(ts, us, ar, tmpl)
	authHandler := handler.NewAuthHandler(us, ts, rs, ar, tmpl)

	router.HandleFunc("/", profileHandler.GetUser).Methods("GET")

	router.HandleFunc("/auth/login", authHandler.LogIn).Methods("GET")
	router.HandleFunc("/auth/logout", authHandler.LogOut).Methods("GET")
	router.HandleFunc("/auth/callback", authHandler.Callback).Methods("GET")

	router.HandleFunc("/ds/resumes", profileHandler.GetResumes).Methods("GET")

	log.Println("starting server at 44444")
	log.Println("login: ", "http://localhost:44444/auth/login")
	log.Fatal(http.ListenAndServe(":44444", router))
}
