package mcproxy

import (
	"github.com/go-logr/logr"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func (p *MCProxy) HandlePreLogin(e *proxy.PreLoginEvent) {
	ctx := e.Conn().Context()
	log := logr.FromContextOrDiscard(ctx)
	log = log.WithValues(
		"username", e.Username(),
	)
	ctx = logr.NewContext(ctx, log)

	err := p.servers.CheckUserAuthorized("vanilla", e.Username())
	if err != nil {
		log.Error(err, "Unauthorized login")
		e.Deny(nil)
		return
	}

	err = p.servers.PrepareServer(ctx, p, "vanilla")
	if err != nil {
		log.Error(err, "Problem preparing server")
		e.Deny(nil)
		return
	}

	log.Info("Allowing connection")
	e.Allow()
}

func (p *MCProxy) HandlePlayerChooseInitialServerEvent(e *proxy.PlayerChooseInitialServerEvent) {
	ctx := e.Player().Context()
	log := logr.FromContextOrDiscard(ctx)

	e.SetInitialServer(p.Server("vanilla"))

	log.Info("Hello")
}

func (p *MCProxy) HandlePlayerConnected(e *proxy.ServerPostConnectEvent) {
	ctx := e.Player().Context()
	log := logr.FromContextOrDiscard(ctx)

	player := e.Player()
	server := player.CurrentServer().Server()

	log.WithValues("connected", server.Players())
	log.Info("Player has connected")
}

func (p *MCProxy) HandlePlayerDisconnected(e *proxy.DisconnectEvent) {
	ctx := e.Player().Context()
	log := logr.FromContextOrDiscard(ctx)

	log.Info("Player has disconnected")
}
