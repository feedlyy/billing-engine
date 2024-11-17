package main

import (
	"billingg-engine/handler"
	"billingg-engine/router"
	md "billingg-engine/router/middleware"
	"billingg-engine/service"
	"billingg-engine/util"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	users, loans, paymentHistories := util.GenerateDummyData()
	svc := service.NewLoanService(users, loans, paymentHistories)
	loanHandler := handler.NewLoanHandler(svc)
	middleware := md.NewMiddleware(users)
	srv := router.NewServer()
	router.RegistRoutes(*srv, loanHandler, middleware)

	logrus.Info("starting application at port 3000")
	if err := srv.Start(":3000"); err != nil || !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
