package manager

import (
	"context"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type ServerManager interface {
	CheckUserAuthorized(string, string) error

	StartServer(context.Context, proxy.ServerRegistry, string) error
}
