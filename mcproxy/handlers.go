package mcproxy

import (
	"github.com/go-logr/logr"
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
	)
	ctx = logr.NewContext(ctx, log)

	err := p.servers.CheckUserAuthorized(VANILLA, e.Username())
	if err != nil {
		log.Error(err, "Unauthorized login")
		e.Deny(nil)
		return
	}

	err = p.servers.StartServer(ctx, p, VANILLA)
	if err != nil {
		log.Error(err, "Problem preparing server")
		e.Deny(nil)
		return
	}

	log.Info("Allowing connection")
	e.Allow()
}

func (p *MCProxy) HandleLogin(e *proxy.LoginEvent) {
	ctx := e.Player().Context()
	log := logr.FromContextOrDiscard(ctx)

	log.WithValues(
		"player", e.Player(),
	).Info("LoginEvent")
	e.Allow()
}

// HandlePlayerChooseInitialServer quickly selects the server
// for the user (expects server to be in registry, else fails)
func (p *MCProxy) HandlePlayerChooseInitialServer(e *proxy.PlayerChooseInitialServerEvent) {
	ctx := e.Player().Context()
	log := logr.FromContextOrDiscard(ctx)

	log.WithValues(
		"vanilla_server", p.Server(VANILLA),
	).Info("ARE WE CHOOSING?")
	e.SetInitialServer(p.Server(VANILLA))

	log.Info("Hello")
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
	ctx := e.Player().Context()
	log := logr.FromContextOrDiscard(ctx)

	log = log.WithValues(
		"server_addr", e.Server().ServerInfo().Addr().String(),
	)
	log.Info("Kicked JOSEPH")
}

// HandlePlayerDisconnected is called when the player disconnects from
// the Proxy. Keep track of how many players are connected, shut down
// server if count is 0.
func (p *MCProxy) HandlePlayerDisconnected(e *proxy.DisconnectEvent) {
	ctx := e.Player().Context()
	log := logr.FromContextOrDiscard(ctx)

	log = log.WithValues(
		"username", e.Player().Username(),
		"host", e.Player().VirtualHost().String(),
	)
	log.Info("Player has disconnected")
}
