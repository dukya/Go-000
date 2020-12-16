package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"geektime/Go-000/Week04/api/user/v1"
	"geektime/Go-000/Week04/internal/biz"
	"geektime/Go-000/Week04/internal/data"
)

type UserService struct {
	v1.UnimplementedUserServer
	uc *biz.UserUseCase
}

func NewUserService(uc *biz.UserUseCase) *UserService {
	return &UserService{uc: uc}
}

func (srv *UserService) GetUserById(ctx context.Context, r *v1.GetUserByIdRequest) (*v1.GetUserByIdReply, error) {
	// TODO: dto对象应该变为do对象
	doId := r.Id

	user, err := srv.uc.GetUserInfoById(ctx, doId)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user is not found.")
		}
		return nil, status.Errorf(codes.Internal, "error:%v", err)
	}
	return &v1.GetUserByIdReply{Id: user.Id, Name: user.Name, Mobile: user.Mobile}, nil
}
