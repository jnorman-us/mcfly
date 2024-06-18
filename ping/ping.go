package ping

import (
	"context"
	"time"

	mcpinger "github.com/Raqbit/mc-pinger"
	"github.com/go-logr/logr"
	"github.com/jnorman-us/mcfly/mcserver/server"
)

const PingResponseDuration = 1 * time.Second
const WaitForServerDuration = 3 * time.Minute

func PingServer(ctx context.Context, s server.Server) (*mcpinger.ServerInfo, error) {
	ctx, _ = context.WithTimeout(ctx, PingResponseDuration)
	pinger := mcpinger.New(
		s.PrivateHost, 25565,
		mcpinger.WithContext(ctx),
	)
	return pinger.Ping()
}

func WaitForServerStatus(ctx context.Context, s server.Server, timeout time.Duration) (*mcpinger.ServerInfo, error) {
	log := logr.FromContextOrDiscard(ctx)

	ticker := time.NewTicker(PingResponseDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.V(1).Info("Timed out waiting for server!")
			return nil, ctx.Err()
		case <-ticker.C:
			info, err := PingServer(ctx, s)
			if err != nil {
				log.V(1).Info("No ping response, retrying...")
				continue
			}
			log.V(1).Info("Received ping response!")
			return info, nil
		}
	}
}
