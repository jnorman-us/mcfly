package config

import (
	"github.com/jnorman-us/mcfly/fly/wirefmt"
)

type ServerConfig struct {
	Name      string
	Whitelist []string

	wirefmt.CPUKind
	MemoryMB  int
	StorageGB int

	Image string
}
