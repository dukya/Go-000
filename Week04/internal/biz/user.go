package biz

import (
	"context"
)

//model层的结构定义
type User struct {
	Id     int64
	Name   string
	Mobile string
}

type UserRepo interface {
	GetUserById(id int64) (*User, error)
}

func NewUserUseCase(repo UserRepo) *UserUseCase {
	return &UserUseCase{repo: repo}
}

type UserUseCase struct {
	repo UserRepo
}

func (uc *UserUseCase) GetUserInfoById(ctx context.Context, id int64) (*User, error) {
	// TODO: 进行相关逻辑处理

	return uc.repo.GetUserById(id)
}
