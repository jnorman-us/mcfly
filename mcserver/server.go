package mcserver

import (
	"fmt"

	"github.com/jnorman-us/mcfly/fly/wirefmt"
	"github.com/jnorman-us/mcfly/mcserver/config"
)

type Server config.ServerConfig

func (s Server) Host() string {
	return fmt.Sprintf("%s.vm.mcfly.internal", s.Name)
}

func (s Server) CreateVolumeInput() wirefmt.CreateVolumeInput {
	return wirefmt.CreateVolumeInput{
		Name:   s.Name,
		SizeGB: s.StorageGB,
	}
}
