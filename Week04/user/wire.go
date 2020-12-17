//+build wireinject

package main

import (
	"context"
	"net/http"

	"github.com/Xia-Jialin/Go-000/Week04/user/dao"
	"github.com/Xia-Jialin/Go-000/Week04/user/service"
	"github.com/Xia-Jialin/Go-000/Week04/user/transport"
	"github.com/google/wire"
	"github.com/gorilla/mux"
)

func InitHttpHandler(userDAO dao.UserDAO, ctx context.Context) http.Handler {
	wire.Build(service.MakeUserServiceImpl, NewUserEndpoints, transport.MakeHttpHandler)
	return &mux.Router{}
}
