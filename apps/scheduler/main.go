package main

import (
	"context"
	"log"
	"os"
	"pkg/database"
	"pkg/headhunter"
	"pkg/model"
	"pkg/repository"
	"pkg/service"
	"time"

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
		var errs string
		ctx := context.Background()
		client := headhunter.NewHHClient(ctx, &row.AccessToken, &row.RefreshToken)
		_, err := client.GetUser(ctx)
		if err == headhunter.ErrHHTokenExpired {
			token, err := client.RefreshToken(ctx)
			if err != nil {
				errs += err.Error()
			}
			if err := service.CreateOrUpdateToken(ctx, &model.Token{
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
			}, row.UserID); err != nil {
				errs += err.Error()
			}
			client.AT = &token.AccessToken
			client.RT = &token.RefreshToken
		}
		err = client.BumpResume(ctx, row.ResumeID)
		if err != nil {
			errs += err.Error()
		}
		service.SaveResult(row, time.Now().Format(time.RFC3339), errs)
	}
}
