package cloud

import (
	"context"

	"github.com/jnorman-us/mcfly/fly/wirefmt"
)

type CloudClient interface {
	ListMachines(context.Context) (wirefmt.ListMachinesOutput, error)
	GetMachine(context.Context, string) (*wirefmt.GetMachineOutput, error)
	CreateMachine(context.Context, wirefmt.CreateMachineInput) (*wirefmt.CreateMachineOutput, error)
	StartMachine(context.Context, string) error
	StopMachine(context.Context, string) error

	ListVolumes(context.Context) (wirefmt.ListVolumesOutput, error)
	CreateVolume(context.Context, wirefmt.CreateVolumeInput) (*wirefmt.CreateVolumeOutput, error)
}
