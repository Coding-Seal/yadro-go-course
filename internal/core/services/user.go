package services

import "yadro-go-course/internal/core/ports"

type UserService struct {
	ports.UserRepo
}

func NewUserService(repo ports.UserRepo) *UserService {
	return &UserService{repo}
}
