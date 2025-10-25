package main

import (
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
	"github.com/rs/cors"
)

func main() {
	log.Println("test")
	dbPath := os.Getenv("DB_PATH")

	if err := os.Remove(dbPath); err != nil {
		log.Fatal(err)
	}

	db, err := database.NewSqliteDatabase(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	repository := repository.NewSqliteRepository(db)
	service := service.NewSqliteService(repository)
	auth := auth.NewAuthRepository()

	router := mux.NewRouter()
	authHandler := handler.NewAuthHandler(service, auth)
	profileHandler := handler.NewProfileHandler(service, auth)

	router.HandleFunc("/profile", profileHandler.Profile).Methods("GET")
	router.HandleFunc("/me", authHandler.Me).Methods("GET")

	router.HandleFunc("/auth/login", authHandler.LogIn).Methods("GET")
	router.HandleFunc("/auth/logout", authHandler.LogOut).Methods("GET")
	router.HandleFunc("/auth/callback", authHandler.Callback).Methods("GET")

	router.HandleFunc("/resumes", profileHandler.GetResumes).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	})
	handlerWithCORS := c.Handler(router)

	log.Println("starting server at 44444")
	log.Println("login: ", "http://localhost:44444/auth/login")
	log.Fatal(http.ListenAndServe(":44444", handlerWithCORS))
}
