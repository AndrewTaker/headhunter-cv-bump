package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"pkg/database"
	"pkg/handler"
	"pkg/repository"
	"pkg/service"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

func main() {
	slogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	slog.SetDefault(slog.New(slogHandler))

	dbPath := os.Getenv("DB_PATH")

	// if err := os.Remove(dbPath); err != nil {
	// 	log.Fatal(err)
	// }

	db, err := database.NewSqliteDatabase(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	repository := repository.NewSqliteRepository(db)
	service := service.NewSqliteService(repository)

	router := mux.NewRouter()
	authHandler := handler.NewAuthHandler(service)
	profileHandler := handler.NewProfileHandler(service)

	router.HandleFunc("/me", profileHandler.Me).Methods("GET")
	router.HandleFunc("/resumes", profileHandler.Resumes).Methods("GET")

	router.HandleFunc("/auth/login", authHandler.LogIn).Methods("GET")
	// router.HandleFunc("/auth/logout", authHandler.LogOut).Methods("GET")
	router.HandleFunc("/auth/callback", authHandler.Callback).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false,
	})
	handlerWithCORS := c.Handler(router)
	loggedMux := handler.LogRequestMiddleware(handlerWithCORS)

	slog.Info("starting server at 44444")
	log.Fatal(http.ListenAndServe(":44444", loggedMux))
}
