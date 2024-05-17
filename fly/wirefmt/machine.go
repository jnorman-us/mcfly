package wirefmt

type ListMachinesOutput []Machine

type CreateMachineInput struct {
	MachineConfig `json:"config"`

	Name   string `json:"name"`
	Region string `json:"region"`

	SkipLaunch bool `json:"skip_launch"`
}

type CreateMachineOutput Machine
type GetMachineOutput Machine

type Machine struct {
	MachineConfig `json:"config"`

	ID        string `json:"id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	PrivateIP string `json:"private_ip"`
}

const MachineStateStarted = "started"
const MachineStateStopped = "stopped"
const MachineStateDestroyed = "destroyed"

type MachineConfig struct {
	Image   string  `json:"image"`
	Guest   Guest   `json:"guest"`
	Mounts  []Mount `json:"mounts"`
	Restart Restart `json:"restart"`

	Env map[string]string `json:"env"`
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

const RestartPolicyNo = "no"
const RestartPolicyAlways = "always"
const RestartPolicyOnFailure = "on-failure"

type Restart struct {
	MaxRetries int    `json:"max_retries"`
	Policy     string `json:"policy"`
}
