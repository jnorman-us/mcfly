package cloud

import (
	"context"

	"github.com/jnorman-us/mcfly/fly/wirefmt"
)

type CloudClient interface {
	CreateMachine(context.Context, wirefmt.CreateMachineInput) (*wirefmt.CreateMachineOutput, error)

	ListVolumes(context.Context) (wirefmt.ListVolumesOutput, error)
	CreateVolume(context.Context, wirefmt.CreateVolumeInput) (*wirefmt.CreateVolumeOutput, error)
}
