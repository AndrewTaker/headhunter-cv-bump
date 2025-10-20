package service

import (
	"pkg/model"
	"pkg/repository"
)

type UserService interface {
	CreateOrUpdateUser(user *model.User) error
	GetUser(id string) (*model.User, error)
	DeleteUser(id string) error
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

func (us *UserServiceImpl) GetUser(id string) (*model.User, error) {
	return us.userRepo.GetUserByID(id)
}

func (us *UserServiceImpl) DeleteUser(id string) error {
	return us.userRepo.DeleteUserByID(id)
}
