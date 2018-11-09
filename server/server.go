package server

import (
	"net"
	"net/http"

	"github.com/pkg/errors"
)

type Server struct {
	listener net.Listener
}

func NewServer(address string) (server *Server, err error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to listen on tcp addr %s", address)
		return
	}

	http.HandleFunc("/profile", HandleProfile)

	server = &Server{
		listener: listener,
	}

	return
}

func (s *Server) Listen() (err error) {
	err = http.Serve(s.listener, nil)
	return
}

func (s *Server) Close() (err error) {
	if s.listener == nil {
		return
	}

	err = s.listener.Close()
	return
}
