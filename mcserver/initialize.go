package mcserver

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/jnorman-us/mcfly/fly/wirefmt"
	"github.com/jnorman-us/mcfly/mcserver/manager"
)

func (m *CloudServerManager) Initialize(ctx context.Context) {
	log := logr.FromContextOrDiscard(ctx)
	log.Info("Initializing servers")

	log.V(1).Info("Collecting existing infrastructure")
	volumesList, err := m.cloud.ListVolumes(ctx)
	if err != nil {
		log.Error(err, "Failed to list volumes")
		return
	}
	machinesList, err := m.cloud.ListMachines(ctx)
	if err != nil {
		log.Error(err, "Failed to list machines")
		return
	}
	existVols := filterVolumes(volumesList)
	existMachines := filterMachines(machinesList)
	log.V(1).WithValues(
		"volumes", existVols,
		"machines", existMachines,
	).Info("Retrieved existing infrastructure")

	for _, server := range m.servers {
		if err := m.prepareServer(ctx, server, existVols, existMachines); err != nil {
			log.Error(err, "Failed to prepare server", "server", server.Name)
		}
	}
}

func (m *CloudServerManager) prepareServer(
	ctx context.Context,
	s *Server,
	volumes map[string]wirefmt.Volume,
	machines map[string]wirefmt.Machine,
) error {
	log := logr.FromContextOrDiscard(ctx).WithValues("server", s.Name)

	log.V(1).Info("Preparing server", "server", s.Name)

	if volume, ok := volumes[s.Name]; ok {
		log.V(1).Info("Found existing volume")
		s.VolumeID = volume.ID
	} else {
		volume, err := m.prepareVolume(ctx, s)
		if err != nil {
			log.Error(err, "Failed to prepare volume")
			return fmt.Errorf("%w: %w", manager.ErrorCloud, err)
		}
		s.VolumeID = volume.ID
	}

	if machine, ok := machines[s.Name]; ok {
		log.V(1).Info("Found existing machine")
		s.MachineID = machine.ID
	} else {
		machine, err := m.prepareMachine(ctx, s)
		if err != nil {
			log.Error(err, "Failed to prepare machine")
			return fmt.Errorf("%w: %w", manager.ErrorCloud, err)
		}
		s.MachineID = machine.ID
	}

	return nil
}

func (m *CloudServerManager) prepareVolume(ctx context.Context, s *Server) (*wirefmt.Volume, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("server", s.Name)

	input := s.CreateVolumeInput()
	log.Info("Creating volume in cloud...", "input", input)

	output, err := m.cloud.CreateVolume(ctx, input)
	if err != nil {
		return nil, err
	}

	volume := wirefmt.Volume(*output)
	return &volume, nil
}

func (m *CloudServerManager) prepareMachine(ctx context.Context, s *Server) (*wirefmt.Machine, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("server", s.Name)

	// create nonexistent machine
	input := s.CreateMachineInput()
	log.Info("Creating machine in cloud...", "input", input)

	output, err := m.cloud.CreateMachine(ctx, input)
	if err != nil {
		return nil, err
	}

	machine := wirefmt.Machine(*output)
	return &machine, nil
}

func filterVolumes(o wirefmt.ListVolumesOutput) map[string]wirefmt.Volume {
	var volumes = map[string]wirefmt.Volume{}
	for _, v := range o {
		if v.State != wirefmt.VolumeStateCreated {
			continue
		}
		volumes[v.Name] = v
	}
	return volumes
}

func filterMachines(o wirefmt.ListMachinesOutput) map[string]wirefmt.Machine {
	var machines = map[string]wirefmt.Machine{}
	for _, m := range o {
		machines[m.Name] = m
	}
	return machines
}
