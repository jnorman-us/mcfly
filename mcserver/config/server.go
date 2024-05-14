package config

import "github.com/jnorman-us/mcfly/fly/wirefmt"

type ServerConfig struct {
	Name      string
	Whitelist []string

	CPUKind   string
	CPUs      int
	MemoryMB  int
	StorageGB int

	Image string

	Restart wirefmt.Restart

	Env map[string]string
}
