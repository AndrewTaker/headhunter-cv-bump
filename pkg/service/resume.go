package service

import (
	"pkg/model"
	"pkg/repository"
)

type ResumeService interface {
	CreateOrUpdateResumes(resumes []model.Resume, userID string) error
	GetUserResumes(userID string) ([]model.Resume, error)
	GetUserResume(resumeID, userID string) (*model.Resume, error)
	DeleteUserResumes(resumes []model.Resume, userID string) error
	ToggleResumeScheduling(resumeID, userID string, isScheduled bool) error
}

type ResumeServiceImpl struct {
	resumeRepo repository.ResumeRepository
}

func NewResumeService(rr repository.ResumeRepository) ResumeService {
	return &ResumeServiceImpl{resumeRepo: rr}
}

func (rs *ResumeServiceImpl) CreateOrUpdateResumes(resumes []model.Resume, userID string) error {
	return rs.resumeRepo.CreateOrUpdateResumes(resumes, userID)
}

func (rs *ResumeServiceImpl) GetUserResumes(userID string) ([]model.Resume, error) {
	return rs.resumeRepo.GetUserResumes(userID)
}

func (rs *ResumeServiceImpl) GetUserResume(resumeID, userID string) (*model.Resume, error) {
	return rs.resumeRepo.GetResumeByID(resumeID, userID)
}

func (rs *ResumeServiceImpl) ToggleResumeScheduling(resumeID, userID string, isScheduled bool) error {
	return rs.resumeRepo.ToggleScheduling(resumeID, userID, isScheduled)
}

func (rs *ResumeServiceImpl) DeleteUserResumes(resumes []model.Resume, userID string) error {
	return rs.resumeRepo.DeleteResumesByUserID(resumes, userID)
}
