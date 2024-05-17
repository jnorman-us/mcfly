package manager

import (
	"context"
)

type ServerManager interface {
	CheckUserAuthorized(string, string) error
	CheckServerReady(context.Context, string) error
	CheckServerStarted(context.Context, string) error

	StartServer(context.Context, string) error
	WaitForServer(context.Context, string) error
	StopServer(context.Context, string) error

	MarkServerReady(string)
	MarkServerHalted(string)
}
