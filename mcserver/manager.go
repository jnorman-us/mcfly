package mcserver

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/jnorman-us/mcfly/fly/cloud"
	"github.com/jnorman-us/mcfly/fly/wirefmt"
	"github.com/jnorman-us/mcfly/mcserver/manager"
	"github.com/jnorman-us/mcfly/mcserver/server"
	"github.com/jnorman-us/mcfly/ping"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var default_servers = map[string]*server.Server{
	"vanilla": {
		ServerConfig: server.ServerConfig{
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
	servers  map[string]*server.Server
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

func (m *CloudServerManager) CheckServerReady(ctx context.Context, name string) error {
	server := m.servers[name]

	_, err := ping.PingServer(ctx, *server)
	if err != nil {
		return fmt.Errorf("%w: %w", manager.ErrorServerNotReady, err)
	}

	return nil
}

func (m *CloudServerManager) CheckServerStarted(ctx context.Context, name string) error {
	server := m.servers[name]
	log := logr.FromContextOrDiscard(ctx).WithValues(
		"machine_id", server.MachineID,
	)

	log.V(1).Info("Getting machine status")
	machine, err := m.cloud.GetMachine(ctx, server.MachineID)
	if err != nil {
		return fmt.Errorf("%w: %w", manager.ErrorCloud, err)
	}

	log.V(1).WithValues(
		"state", machine.State,
	).Info("Got machine status")
	if machine.State != wirefmt.MachineStateStarted {
		return manager.ErrorServerNotStarted
	}
	return nil
}

func (m *CloudServerManager) StartServer(ctx context.Context, name string) error {
	log := logr.FromContextOrDiscard(ctx)
	server := m.servers[name]

	log.V(1).WithValues(
		"machine_id", server.MachineID,
	).Info("Starting machine")

	err := m.cloud.StartMachine(ctx, server.MachineID)
	if err != nil {
		return fmt.Errorf("%w: %w", manager.ErrorCloud, err)
	}
	return nil
}

func (m *CloudServerManager) WaitForServer(ctx context.Context, name string) error {
	server := m.servers[name]

	_, err := ping.WaitForServerStatus(ctx, *server, ping.WaitForServerDuration)
	if err != nil {
		return fmt.Errorf("%w: %w", manager.ErrorTimeout, err)
	}

	return nil
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

func (m *CloudServerManager) MarkServerReady(name string) {
	m.registry.Register(m.servers[name])
}

func (m *CloudServerManager) MarkServerHalted(name string) {
	m.registry.Unregister(m.servers[name])
}
