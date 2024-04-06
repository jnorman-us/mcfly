package mcserver

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/jnorman-us/mcfly/fly/cloud"
	"github.com/jnorman-us/mcfly/fly/wirefmt"
	"github.com/jnorman-us/mcfly/mcserver/config"
	"github.com/jnorman-us/mcfly/mcserver/manager"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var default_servers = map[string]Server{
	"vanilla": {
		ServerConfig: config.ServerConfig{
			Name: "vanilla",
			Whitelist: []string{
				"Wine_Craft",
			},
			CPUKind:   wirefmt.CPUKindShared,
			CPUs:      1,
			MemoryMB:  1024,
			StorageGB: 5,

			Image: "itzg/minecraft-server:latest",
		},
	},
}

type CloudServerManager struct {
	servers map[string]Server
	cloud   cloud.CloudClient
}

func NewCloudServerManager(cc cloud.CloudClient) *CloudServerManager {
	return &CloudServerManager{
		servers: default_servers,
		cloud:   cc,
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

func (m *CloudServerManager) StartServer(ctx context.Context, registry proxy.ServerRegistry, name string) error {
	log := logr.FromContextOrDiscard(ctx).WithValues(
		"server", name,
	)

	log.Info("Starting server!")

	return nil
}

func (m *CloudServerManager) VerifyServer(ctx context.Context, serverName string) error {
	return nil
}
