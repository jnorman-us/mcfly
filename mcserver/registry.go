package mcserver

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/jnorman-us/mcfly/fly/cloud"
	"github.com/jnorman-us/mcfly/fly/wirefmt"
	"github.com/jnorman-us/mcfly/mcserver/manager"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var default_servers = map[string]Server{
	"vanilla": Server{
		Name: "vanilla",
		Whitelist: []string{
			"Wine_Craft",
		},
		CPUKind:   wirefmt.SharedCPU1x,
		MemoryMB:  1024,
		StorageGB: 5,

		Image: "itzg/minecraft-server:latest",
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

func (m *CloudServerManager) PrepareServer(ctx context.Context, registry proxy.ServerRegistry, name string) error {
	log := logr.FromContextOrDiscard(ctx).WithValues(
		"server", name,
	)

	var s, _ = m.servers[name]

	volume, err := m.prepareVolume(ctx, s)
	if err != nil {
		return fmt.Errorf("%w: %w", manager.ErrorCloud, err)
	}

	log.V(1).Info("Proceeding with volume", "volume", volume)

	return nil
}

func (m *CloudServerManager) prepareVolume(ctx context.Context, s Server) (*wirefmt.Volume, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues(
		"server", s.Name,
	)

	// verify volume exists
	volumes, err := m.cloud.ListVolumes(ctx)
	if err != nil {
		return nil, err
	}
	log.V(1).Info("Listed volumes in cloud", "volumes", volumes)
	for _, v := range volumes {
		// filter out unusable volumes
		if v.State != wirefmt.VolumeStateCreated {
			continue
		}

		if v.Name == s.Name {
			return &v, nil
		}
	}

	// create nonexistent volume
	input := s.CreateVolumeInput()
	log.V(1).Info("Creating volume in cloud", "input", input)

	output, err := m.cloud.CreateVolume(ctx, input)
	if err != nil {
		return nil, err
	}

	volume := wirefmt.Volume(*output)
	return &volume, nil
}

func (m *CloudServerManager) VerifyServer(ctx context.Context, serverName string) error {
	return nil
}
