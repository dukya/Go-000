package service

import (
	"geektime/Go-000/Week02/dao"
	"geektime/Go-000/Week02/model"
)

type UserService struct {
}

func (u *UserService) GetUser(name string) (*model.User, error) {
	return dao.GetUserInfo(name)
}
