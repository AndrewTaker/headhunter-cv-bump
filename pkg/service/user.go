package service

import (
	"pkg/model"
	"pkg/repository"
)

type UserService interface {
	CreateOrUpdateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	DeleteUserByID(id string) error
}

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(ur repository.UserRepository) UserService {
	return &UserServiceImpl{userRepo: ur}
}

func (us *UserServiceImpl) CreateOrUpdateUser(user *model.User) error {
	return us.userRepo.CreateOrUpdateUser(user)
}

func (us *UserServiceImpl) GetUserByID(id string) (*model.User, error) {
	return us.userRepo.GetUserByID(id)
}

func (us *UserServiceImpl) DeleteUserByID(id string) error {
	return us.userRepo.DeleteUserByID(id)
}
