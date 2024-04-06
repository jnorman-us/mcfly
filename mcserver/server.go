package mcserver

import (
	"fmt"

	"github.com/jnorman-us/mcfly/fly/wirefmt"
	"github.com/jnorman-us/mcfly/mcserver/config"
)

type Server struct {
	config.ServerConfig
	VolumeID  string
	MachineID string
}

func (s Server) Host() string {
	return fmt.Sprintf("%s.vm.mcfly.internal", s.Name)
}

func (s Server) CreateVolumeInput() wirefmt.CreateVolumeInput {
	return wirefmt.CreateVolumeInput{
		Name:   s.Name,
		SizeGB: s.StorageGB,
	}
}

func (s Server) CreateMachineInput() wirefmt.CreateMachineInput {
	return wirefmt.CreateMachineInput{
		Name: s.Name,
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
		},
	}
}
