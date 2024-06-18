package mcserver

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
)

func (m *CloudServerManager) FindCloudResources(ctx context.Context) error {
	log := logr.FromContextOrDiscard(ctx)

	log.V(1).Info("Collecting existing infrastructure")
	volumesList, err := m.cloud.ListVolumes(ctx)
	if err != nil {
		return fmt.Errorf("failed to list volumes: %w", err)
	}
	machinesList, err := m.cloud.ListMachines(ctx)
	if err != nil {
		return fmt.Errorf("failed to list machines: %w", err)
	}
	existVols := filterVolumes(volumesList)
	existMachines := filterMachines(machinesList)
	log.V(1).WithValues(
		"volumes", existVols,
		"machines", existMachines,
	).Info("Retrieved existing infrastructure")

	for _, server := range m.servers {
		name := server.Name()

		vol, ok := existVols[name]
		if !ok {
			return fmt.Errorf("missing volume for %s", name)
		}
		server.VolumeID = vol.ID

		machine, ok := existMachines[name]
		if !ok {
			return fmt.Errorf("missing machine for %s", name)
		}
		server.MachineID = machine.ID
		err = server.SetAddr(machine.PrivateIP)
		if err != nil {
			return fmt.Errorf("failed to parse private address for %s: %w", machine.ID, err)
		}
	}

	log.V(1).WithValues(
		"servers", m.servers,
	).Info("Infrastructure found for servers")
	return nil
}

func (m *CloudServerManager) PopulateRegistry() {
	for _, server := range m.servers {
		m.registry.Register(server)
	}
}
