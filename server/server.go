package server

import "net/http"

type Server struct {
	server http.Server
}

func SeverConnectionBuilder(hand http.Handler) error{
	serv := &Server{
		server: http.Server{
			Addr: ":8080",
			Handler: hand,
		},
	}

	return serv.server.ListenAndServe()
}