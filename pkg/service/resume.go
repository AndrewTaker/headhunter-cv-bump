package service

import "pkg/model"

type ResumeService interface {
	CreateOrUpdateResume(resumes []model.Resume, userID string) error
	GetResumes(userID string) ([]model.Resume, error)
	GetResume(resumeID, userID string) (*model.Resume, error)
	ToggleScheduling(resumeID, userID string, isScheduled bool) error
}
