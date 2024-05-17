package mcserver

import (
	"context"
	"time"

	mcpinger "github.com/Raqbit/mc-pinger"
)

func PingServer(ctx context.Context, s *Server, timeout time.Duration) (*mcpinger.ServerInfo, error) {
	ctx, _ = context.WithTimeout(ctx, timeout)
	pinger := mcpinger.New(
		s.PrivateHost, 25565,
		mcpinger.WithContext(ctx),
	)
	return pinger.Ping()
}
