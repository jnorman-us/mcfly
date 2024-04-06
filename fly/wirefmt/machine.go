package wirefmt

type ListMachinesOutput []Machine

type CreateMachineInput struct {
	MachineConfig `json:"config"`

	Name   string `json:"name"`
	Region string `json:"region"`
}

type CreateMachineOutput Machine

type Machine struct {
	MachineConfig `json:"config"`

	ID   string `json:"id"`
	Name string `json:"name"`
}

type MachineConfig struct {
	Image  string  `json:"image"`
	Guest  Guest   `json:"guest"`
	Mounts []Mount `json:"mounts"`
}

type Guest struct {
	CPUKind  string `json:"cpu_kind"`
	CPUs     int    `json:"cpus"`
	MemoryMB int    `json:"memory_mb"`
}

type Mount struct {
	Name   string `json:"name"`
	Volume string `json:"volume"`
	Path   string `json:"path"`
}
