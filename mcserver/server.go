package mcserver

import (
	"fmt"

	"github.com/jnorman-us/mcfly/mcserver/config"
)

type Server struct {
	config.ServerConfig
}

func NewServer(cfg config.ServerConfig) Server {
	return Server{
		ServerConfig: cfg,
	}
}

func (s *Server) Host() string {
	return fmt.Sprintf("%s.vm.mcfly.internal", s.Name)
}
