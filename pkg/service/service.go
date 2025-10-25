package service

import (
	"context"
	"database/sql"
	"pkg/model"
	"pkg/repository"
	"time"
)

type SqliteService struct {
	repo *repository.SqliteRepository
}

func NewSqliteService(r *repository.SqliteRepository) *SqliteService {
	return &SqliteService{repo: r}
}

func (rs *SqliteService) CreateOrUpdateResumes(resumes []model.Resume, userID string) error {
	return rs.repo.ResumeCreateOrUpdateBatch(resumes, userID)
}

func (rs *SqliteService) GetUserResumes(userID string) ([]model.Resume, error) {
	return rs.repo.ResumeGetByUserIDBatch(userID)
}

func (rs *SqliteService) GetUserResume(resumeID, userID string) (*model.Resume, error) {
	return rs.repo.ResumeGetByID(resumeID, userID)
}

func (rs *SqliteService) ToggleResumeScheduling(resumeID, userID string, isScheduled bool) error {
	return rs.repo.ResumeToggleScheduling(resumeID, userID, isScheduled)
}

func (rs *SqliteService) DeleteUserResumes(resumes []model.Resume, userID string) error {
	return rs.repo.ResumeDeleteByUserID(resumes, userID)
}

func (rs *SqliteService) GetSchedules() ([]model.JoinedScheduler, error) {
	var schedules []model.JoinedScheduler
	var err error

	if schedules, err = rs.repo.ScheduleGetBatch(); err != nil {
		return schedules, err
	}

	for i := range schedules {
		var token model.Token
		token.AccessToken = schedules[i].AccessToken
		token.RefreshToken = schedules[i].RefreshToken

		if err := token.Decrypt(); err != nil {
			return schedules, err
		}

		schedules[i].AccessToken = token.AccessToken
		schedules[i].RefreshToken = token.RefreshToken
	}

	return schedules, nil
}

func (rs *SqliteService) SaveResult(s model.JoinedScheduler, timestamp, errors string) error {
	return rs.repo.ScheduleSave(s, timestamp, errors)
}

func (rs *SqliteService) CreateOrUpdateToken(ctx context.Context, token *model.Token, userID string) error {
	if err := token.Encrypt(); err != nil {
		return err
	}

	return rs.repo.TokenSaveOrCreate(ctx, token, userID)
}

func (rs *SqliteService) GetTokenByUserID(ctx context.Context, userID string) (*model.Token, error) {
	token, err := rs.repo.TokenGetByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := token.Decrypt(); err != nil {
		return nil, err
	}

	return token, nil
}

func (rs *SqliteService) CreateOrUpdateUser(user *model.User) error {
	return rs.repo.UserCreateOrUpdate(user)
}

func (rs *SqliteService) GetUser(id string) (*model.User, error) {
	return rs.repo.UserGetByID(id)
}

func (rs *SqliteService) DeleteUser(id string) error {
	return rs.repo.UserDeleteByID(id)
}

func (rs *SqliteService) SaveSession(ctx context.Context, sessID, userID string, expiresAt time.Time) error {
	return rs.repo.SessionSave(ctx, sessID, userID, expiresAt)
}

func (rs *SqliteService) DeleteSession(ctx context.Context, sessID, userID string) error {
	return rs.repo.SessionDelete(ctx, sessID, userID)
}

func (rs *SqliteService) GetUserBySession(ctx context.Context, sessID string) (*model.User, error) {
	return rs.repo.UserGetBySessionID(ctx, sessID)
}
