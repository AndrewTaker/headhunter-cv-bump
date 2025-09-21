package main

import (
	"log"
	"pkg/database"
	"pkg/repository"
	"pkg/service"

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
