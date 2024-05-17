package manager

import (
	"context"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type ServerManager interface {
	CheckUserAuthorized(string, string) error
	GetRunningServer(context.Context, string) (proxy.RegisteredServer, error)

	StartServer(context.Context, string) error
	StopServer(context.Context, string) error
	MarkServerHalted(context.Context, string)
}
