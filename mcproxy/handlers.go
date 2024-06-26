package mcproxy

import (
	"errors"

	"github.com/go-logr/logr"
	"github.com/jnorman-us/mcfly/mcserver/manager"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

const VANILLA = "vanilla"

// HandlePreLogin is called when a player connects to the
// proxy. We can take as much time as needed to spin up
// the server if needed.
func (p *MCProxy) HandlePreLogin(e *proxy.PreLoginEvent) {
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
		p.servers.MarkServerReady(VANILLA)
		e.Allow()
		return
	}

	err = p.servers.CheckServerStarted(ctx, VANILLA)
	log.WithValues("TEST_ERR", err).Info("HELLO")
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
	p.servers.MarkServerReady(VANILLA)
	e.Allow()
}

// HandlePlayerChooseInitialServer quickly selects the server
// for the user (expects server to be in registry, else fails)
func (p *MCProxy) HandlePlayerChooseInitialServer(e *proxy.PlayerChooseInitialServerEvent) {
	e.SetInitialServer(p.Server(VANILLA))
}

// HandlePlayerConnected gets called when a player finishes the login
// process and the server gets started. Nothing needs to be done here...
func (p *MCProxy) HandlePlayerConnected(e *proxy.ServerPostConnectEvent) {
	ctx := e.Player().Context()
	log := logr.FromContextOrDiscard(ctx)

	player := e.Player()
	server := player.CurrentServer().Server()

	log = log.WithValues(
		"connected_players", server.Players(),
		"server_info", server.ServerInfo(),
	)
	log.Info("Player has connected")
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
	)
	log.Info("Player has disconnected")

	//err := p.servers.StopServer(ctx, VANILLA)
	//if err != nil {
	//	log.Error(err, "Problem stopping server")
	//}
}
