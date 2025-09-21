package service

import (
	"pkg/model"
	"pkg/repository"
)

type SchedulerService interface {
	GetSchedules() ([]model.JoinedScheduler, error)
	SaveResult(s model.JoinedScheduler, timestamp, err string) error
}

type SchedulerServiceImpl struct {
	schedulerRepo repository.SchedulerRepository
}

func NewSchedulerService(ss repository.SchedulerRepository) SchedulerService {
	return &SchedulerServiceImpl{schedulerRepo: ss}
}

func (ss *SchedulerServiceImpl) GetSchedules() ([]model.JoinedScheduler, error) {
	var schedules []model.JoinedScheduler
	var err error

	if schedules, err = ss.schedulerRepo.GetSchedules(); err != nil {
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

func (ss *SchedulerServiceImpl) SaveResult(s model.JoinedScheduler, timestamp, errors string) error {
	return ss.schedulerRepo.SaveResult(s, timestamp, errors)
}
