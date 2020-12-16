// +build wireinject

package main

import (
	"github.com/google/wire"

	"geektime/Go-000/Week04/internal/biz"
	"geektime/Go-000/Week04/internal/data"
	"geektime/Go-000/Week04/internal/pkg/grpctransport"
	"geektime/Go-000/Week04/internal/service"
)

func InitializeServer() (*grpctransport.Server, func(), error) {
	wire.Build(grpctransport.NewServer, service.NewUserService, biz.NewUserUseCase, data.Provider)
	return nil, nil, nil
}
