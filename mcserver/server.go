package mcserver

import (
	"net"

	"github.com/jnorman-us/mcfly/fly/wirefmt"
	"github.com/jnorman-us/mcfly/mcserver/config"
)

type Server struct {
	config.ServerConfig
	VolumeID    string
	MachineID   string
	PrivateHost string
	PrivateAddr net.Addr
}

func (s Server) Name() string {
	return s.ServerConfig.Name
}

func (s Server) Addr() net.Addr {
	return s.PrivateAddr
}

func (s Server) CreateVolumeInput() wirefmt.CreateVolumeInput {
	return wirefmt.CreateVolumeInput{
		Name:   s.Name(),
		SizeGB: s.StorageGB,
	}
}

func (s Server) CreateMachineInput() wirefmt.CreateMachineInput {
	return wirefmt.CreateMachineInput{
		Name: s.Name(),
		MachineConfig: wirefmt.MachineConfig{
			Image: s.Image,
			Guest: wirefmt.Guest{
				CPUKind:  s.CPUKind,
				CPUs:     s.CPUs,
				MemoryMB: s.MemoryMB,
			},
			Mounts: []wirefmt.Mount{{
				Name:   "data",
				Volume: s.VolumeID,
				Path:   "/data",
			}},
			Restart: s.Restart,
			Env:     s.Env,
		},
		SkipLaunch: true,
	}
}
