package server

import (
	"net"
	"net/http"

	"code.cloudfoundry.org/lager"
	"github.com/pkg/errors"
)

type Server struct {
	logger          lager.Logger
	listener        net.Listener
	assetsDirectory string
}

func NewServer(logger lager.Logger, address, assetsDirectory string) (server *Server, err error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to listen on tcp addr %s", address)
		return
	}

	server = &Server{
		logger:          logger,
		listener:        listener,
		assetsDirectory: assetsDirectory,
	}

	http.HandleFunc("/profile", server.HandleProfile)
	http.HandleFunc("/flamegraph", server.HandleFlamegraph)
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir(assetsDirectory))))

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
