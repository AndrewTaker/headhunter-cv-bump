package main

import (
	"log"
	"net/http"
	"pkg/auth"
	"pkg/database"
	"pkg/handler"
	"pkg/repository"
	"pkg/service"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.NewSqliteDatabase("/home/pepega/pepehands/hh-cv/cvu-hh.db")
	if err != nil {
		log.Fatal(err)
	}

	ur := repository.NewSqliteUserRepository(db)
	us := service.NewUserService(ur)

	tr := repository.NewSqliteTokenRepository(db)
	ts := service.NewTokenService(tr)

	rr := repository.NewSqliteResumeRepository(db)
	rs := service.NewResumeService(rr)

	sr := repository.NewSqliteSchedulerRepository(db)
	ss := service.NewSchedulerService(sr)

	ar := auth.NewAuthRepository()

	router := mux.NewRouter()
	userHandler := handler.NewUserHandler(us, ar)
	authHandler := handler.NewAuthHandler(us, ts, rs, ar)
	router.HandleFunc("/users", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/auth/login", authHandler.LogIn).Methods("GET")
	router.HandleFunc("/auth/callback", authHandler.Callback).Methods("GET")

	log.Println("starting server at 44444")
	log.Println("login: ", "http://localhost:44444/auth/login")
	log.Fatal(http.ListenAndServe(":44444", router))

	user, err := us.GetUser("60645454")
	if err != nil {
		log.Fatal(err)
	}

	usersToken, err := ts.GetToken(user.ID)
	if err != nil {
		log.Fatal(err)
	}

	usersResume, err := rs.GetUserResumes(user.ID)
	if err != nil {
		log.Fatal(err)
	}

	schedule, err := ss.GetSchedules()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(user)
	log.Println(usersToken)
	for _, r := range usersResume {
		log.Println(r.Title, r.ID)
	}

	log.Println("----------------------------------------------------")
	for _, s := range schedule {
		log.Println(s.AccessToken, s.RefreshToken, s.ResumeID)
	}
}
