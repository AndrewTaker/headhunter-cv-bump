package main

import (
	"context"
	"log"
	"os"
	"pkg/database"
	"pkg/headhunter"
	"pkg/repository"
	"pkg/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.NewSqliteDatabase(os.Getenv("DB_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	repository := repository.NewSqliteRepository(db)
	service := service.NewSqliteService(repository)

	data, err := service.GetSchedules()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range data {
		ctx := context.Background()
		client := headhunter.NewHHClient(ctx, &row.AccessToken, &row.RefreshToken)
		err := client.BumpResume(ctx, row.ResumeID)
		if err != nil {
			log.Println("ERR", err)
		}
	}
}
