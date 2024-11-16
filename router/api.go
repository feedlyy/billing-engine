package router

import "billingg-engine/handler"

func RegistRoutes(srv Server, handler handler.LoanHandler) {
	srv.AddRoute("GET", "/outstanding", handler.GetCurrentOutStanding)
}
