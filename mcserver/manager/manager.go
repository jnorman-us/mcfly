package manager

import (
	"context"
)

type ServerManager interface {
	CheckUserAuthorized(string, string) error
	CheckServerReady(context.Context, string) error
	CheckServerStarted(context.Context, string) error
	CheckServerEmpty(context.Context, string) error

	StartServer(context.Context, string) error
	WaitForServer(context.Context, string) error

	PrepareHaltServer(string) error
	CancelHaltServer(string)
}
