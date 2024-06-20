package mcproxy

import (
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/jnorman-us/mcfly/halter"
	"github.com/jnorman-us/mcfly/mcserver/manager"
	"go.minekube.com/common/minecraft/color"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

const VANILLA = "vanilla"

// HandlePreLogin is called when a player connects to the
// proxy. We can take as much time as needed to spin up
// the server if needed.
func (p *MCProxy) HandlePreLogin(e *proxy.PreLoginEvent) {
	e.Allow()
	return
	ctx := e.Conn().Context()
	log := logr.FromContextOrDiscard(ctx)
	log = log.WithValues(
		"username", e.Username(),
		"server", VANILLA,
	)
	ctx = logr.NewContext(ctx, log)

	err := p.servers.CheckUserAuthorized(VANILLA, e.Username())
	if err != nil {
		log.Error(err, "Unauthorized login")
		e.Deny(nil)
		return
	}

	err = p.servers.CheckServerReady(ctx, VANILLA)
	if err == nil {
		p.servers.CancelHaltServer(VANILLA)
		e.Allow()
		return
	}

	err = p.servers.CheckServerStarted(ctx, VANILLA)
	if errors.Is(err, manager.ErrorCloud) {
		log.Error(err, "Problem checking if server started")
		e.Deny(nil)
		return
	}
	if errors.Is(err, manager.ErrorServerNotStarted) {
		log.Info("Server not running, starting...")
		err = p.servers.StartServer(ctx, VANILLA)
		if err != nil {
			log.Error(err, "Problem starting server")
			e.Deny(nil)
			return
		}
	}

	log.Info("Waiting for server...")
	err = p.servers.WaitForServer(ctx, VANILLA)
	if err != nil {
		log.Error(err, "Server failed to respond in time")
		e.Deny(nil)
		return
	}

	log.Info("Server ready, allowing connection")
	p.servers.CancelHaltServer(VANILLA)
	e.Allow()
}

func (p *MCProxy) HandleLogin(e *proxy.LoginEvent) {
	fmt.Println("Hello")
	e.Deny(&component.Text{
		Content: "An error occurred while running this command.",
		S:       component.Style{Color: color.Red},
	})
}

// HandlePlayerChooseInitialServer quickly selects the server
// for the user (expects server to be in registry, else fails)
func (p *MCProxy) HandlePlayerChooseInitialServer(e *proxy.PlayerChooseInitialServerEvent) {
	e.SetInitialServer(p.Server(VANILLA))
}

// HandlePlayerKicked is called when the underlying server goes down
// and the player is kicked off. Human intervention required..?
// Least we can do is remove the server from the registry, though we
// might need to mark it as Do Not Resuscitate
func (p *MCProxy) HandlePlayerKicked(e *proxy.KickedFromServerEvent) {
	// ctx := e.Player().Context()
	// log := logr.FromContextOrDiscard(ctx)

	// log = log.WithValues(
	// 	"player", e.Player().GameProfile().Name,
	// 	"server", VANILLA,
	// )

	// log.WithValues("e", e.KickedDuringServerConnect(), "t", e.OriginalReason()).Info("Test")
	// if e.KickedDuringServerConnect() {
	// 	log.Info("Player could not connect to upstream, marking as halted...")
	// 	p.servers.MarkServerHalted(ctx, VANILLA)
	// }
}

// HandlePlayerDisconnected is called when the player disconnects from
// the Proxy. Keep track of how many players are connected, shut down
// server if count is 0.
func (p *MCProxy) HandlePlayerDisconnected(e *proxy.DisconnectEvent) {
	ctx := e.Player().Context()
	log := logr.FromContextOrDiscard(ctx)
	log = log.WithValues(
		"username", e.Player().Username(),
		"server", VANILLA,
		"login_status", e.LoginStatus(),
	)
	log.Info("Player has disconnected")

	err := p.servers.CheckServerEmpty(ctx, VANILLA)
	if errors.Is(err, manager.ErrorServerNotEmpty) {
		log.V(1).Info("Server is not empty")
		return
	}

	err = p.servers.PrepareHaltServer(VANILLA)
	if err != nil {
		log.Error(err, "Failed to prepare halting server")
		return
	}
	log.Info("Server is preparing to halt", "wait_time", halter.HaltWaitDuration)
}
