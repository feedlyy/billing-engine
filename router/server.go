package router

import (
	"net/http"
	"strings"
)

type Server struct {
	handler http.Handler
}

func NewServer() *Server {
	return &Server{
		handler: http.NewServeMux(),
	}
}

func (s *Server) AddRoute(method, pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
	s.handler.(*http.ServeMux).HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if strings.EqualFold(method, r.Method) {
			handlerFunc(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.handler)
}
