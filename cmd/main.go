package main

import (
	"billingg-engine/handler"
	"billingg-engine/router"
	"billingg-engine/service"
	"billingg-engine/util"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	users, loans, paymentHistories := util.GenerateDummyData()
	svc := service.NewLoanService(users, loans, paymentHistories)
	handler := handler.NewLoanHandler(svc)
	middleware := router.NewMiddleware(users)
	srv := router.NewServer()
	router.RegistRoutes(*srv, handler, middleware)

	logrus.Info("starting application at port 3000")
	if err := srv.Start(":3000"); err != nil || !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
