package mcserver

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/jnorman-us/mcfly/fly/cloud"
	"github.com/jnorman-us/mcfly/fly/wirefmt"
	"github.com/jnorman-us/mcfly/mcserver/config"
	"github.com/jnorman-us/mcfly/mcserver/manager"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var default_servers = map[string]*Server{
	"vanilla": {
		ServerConfig: config.ServerConfig{
			Name: "vanilla",
			Whitelist: []string{
				"Wine_Craft",
			},
			CPUKind:   wirefmt.CPUKindShared,
			CPUs:      1,
			MemoryMB:  2048,
			StorageGB: 5,

			Image: "itzg/minecraft-server:latest",

			Restart: wirefmt.Restart{
				Policy: wirefmt.RestartPolicyNo,
			},
			Env: map[string]string{
				"EULA":        "TRUE",
				"VERSION":     "1.20.4",
				"ONLINE_MODE": "FALSE",
			},
		},
	},
}

type CloudServerManager struct {
	servers map[string]*Server
	cloud   cloud.CloudClient
}

func NewCloudServerManager(cc cloud.CloudClient) *CloudServerManager {
	return &CloudServerManager{
		servers: default_servers,
		cloud:   cc,
	}
}

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
		name := server.Name

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
	}

	log.V(1).WithValues(
		"servers", m.servers,
	).Info("Infrastructure found for servers")
	return nil
}

func (m *CloudServerManager) CheckUserAuthorized(name string, username string) error {
	if server, ok := m.servers[name]; ok {
		var whitelist = server.Whitelist
		for _, player := range whitelist {
			if username == player {
				return nil
			}
		}
		return manager.ErrorNotAuthorized
	} else {
		return manager.ErrorServerNotRegistered
	}
}

func (m *CloudServerManager) StartServer(ctx context.Context, registry proxy.ServerRegistry, name string) error {
	log := logr.FromContextOrDiscard(ctx).WithValues(
		"server", name,
	)

	server := m.servers["vanilla"]

	host, err := server.Host()
	if err != nil {
		return fmt.Errorf("failed to parse host: %w", err)
	}

	registry.Register(proxy.NewServerInfo(server.Name, host))

	log.Info("Starting server!")

	return nil
}

func (m *CloudServerManager) VerifyServer(ctx context.Context, serverName string) error {
	return nil
}
