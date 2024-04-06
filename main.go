package main

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/jnorman-us/mcfly/env"
	"github.com/jnorman-us/mcfly/fly"
	"github.com/jnorman-us/mcfly/mcproxy"
	"github.com/jnorman-us/mcfly/mcserver"
	jconfig "go.minekube.com/gate/pkg/edition/java/config"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/gate"
	"go.minekube.com/gate/pkg/gate/config"
	"go.uber.org/zap"
)

var default_config = config.DefaultConfig

func main() {
	ctx := context.Background()

	cfg := env.Config{}
	cfg.FlyToken = os.Getenv(env.KeyFlyToken)
	cfg.FlyApp = "mcfly"

	client := fly.NewFlyClient(cfg)
	manager := mcserver.NewCloudServerManager(client)

	zapLogger := zap.Must(zap.NewDevelopment())
	ctx = logr.NewContext(ctx, zapr.NewLogger(zapLogger))

	proxy.Plugins = append(proxy.Plugins, proxy.Plugin{
		Name: "MCFlyProxy",
		Init: func(ctx context.Context, proxy *proxy.Proxy) error {
			return mcproxy.NewMCProxy(proxy, manager).Init(ctx)
		},
	})

	default_config.Editions.Java.Config.Forwarding.Mode = jconfig.NoneForwardingMode
	default_config.Editions.Java.Config.OnlineMode = false

	// default_config.Config.Servers["vanilla"] = "127.0.0.1:25566"

	gate.Start(ctx, gate.WithConfig(default_config))
}
