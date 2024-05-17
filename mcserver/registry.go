package mcserver

import (
	"context"
	"fmt"
	"net"

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
	servers  map[string]*Server
	registry proxy.ServerRegistry
	cloud    cloud.CloudClient
}

func NewCloudServerManager(cc cloud.CloudClient, r proxy.ServerRegistry) *CloudServerManager {
	return &CloudServerManager{
		servers:  default_servers,
		cloud:    cc,
		registry: r,
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
		server.PrivateHost = machine.PrivateIP
		server.PrivateAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("[%s]:25565", machine.PrivateIP))
		if err != nil {
			return fmt.Errorf("failed to parse private address for %s: %w", machine.ID, err)
		}
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

func (m *CloudServerManager) GetRunningServer(ctx context.Context, name string) (proxy.RegisteredServer, error) {
	server := m.servers[name]
	log := logr.FromContextOrDiscard(ctx).WithValues(
		"machine_id", server.MachineID,
	)

	log.V(1).Info("Getting machine status")
	machine, err := m.cloud.GetMachine(ctx, server.MachineID)
	if err != nil {
		return nil, err
	}

	log.V(1).WithValues(
		"state", machine.State,
	).Info("Got machine status")

	if machine.State != wirefmt.MachineStateStarted {
		m.registry.Unregister(server)
		return nil, nil
	}

	registered := m.registry.Server(name)
	if registered != nil {
		log.V(1).Info("Server already registered")
		return registered, nil
	}

	log.V(1).Info("Registering server")
	registered, _ = m.registry.Register(server)
	log.WithValues("registered", registered).Info("Registered server")
	return registered, nil
}

func (m *CloudServerManager) StartServer(ctx context.Context, name string) error {
	log := logr.FromContextOrDiscard(ctx)
	server := m.servers[name]

	log.V(1).WithValues(
		"machine_id", server.MachineID,
	).Info("Starting machine")

	err := m.cloud.StartMachine(ctx, server.MachineID)
	if err != nil {
		return fmt.Errorf("failed to start machine: %w", err)
	}

	host := server.Addr()
	log.V(1).WithValues(
		"host", host,
	).Info("Registering server")

	_, err = m.registry.Register(server)
	return err
}

func (m *CloudServerManager) StopServer(ctx context.Context, name string) error {
	log := logr.FromContextOrDiscard(ctx)
	server := m.servers[name]

	log.V(1).WithValues(
		"machine_id", server.MachineID,
	).Info("Stopping machine")

	err := m.cloud.StopMachine(ctx, server.MachineID)
	if err != nil {
		return fmt.Errorf("failed to stop machine: %w", err)
	}

	m.registry.Unregister(server)
	return nil
}

func (m *CloudServerManager) MarkServerHalted(ctx context.Context, name string) {
	server := m.servers[name]

	m.registry.Unregister(server)
}

func (m *CloudServerManager) VerifyServer(ctx context.Context, serverName string) error {
	return nil
}
