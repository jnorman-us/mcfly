package config

type ServerConfig struct {
	Name      string
	Whitelist []string

	CPUKind   string
	CPUs      int
	MemoryMB  int
	StorageGB int

	Image string
}
