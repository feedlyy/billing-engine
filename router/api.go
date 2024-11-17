package router

import (
	_const "billingg-engine/const"
	"billingg-engine/handler"
	"billingg-engine/router/middleware"
)

func RegistRoutes(srv Server, handler handler.LoanHandler, m middleware.Middleware) {
	srv.AddRoute("POST", "/login", m.Login)
	srv.AddRoute("GET", "/outstanding", middleware.AuthMiddlewareWithRole(handler.GetCurrentOutStanding, _const.RoleCustomer))
	srv.AddRoute("GET", "/check", middleware.AuthMiddlewareWithRole(handler.CheckIsDelinquent, _const.RoleAdmin))
	srv.AddRoute("POST", "/pay", middleware.AuthMiddlewareWithRole(handler.Payment, _const.RoleCustomer))
}
