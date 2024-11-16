package router

import (
	_const "billingg-engine/const"
	"billingg-engine/handler"
)

func RegistRoutes(srv Server, handler handler.LoanHandler, m Middleware) {
	srv.AddRoute("POST", "/login", m.Login)
	srv.AddRoute("GET", "/outstanding", AuthMiddlewareWithRole(handler.GetCurrentOutStanding, _const.RoleCustomer))
}
