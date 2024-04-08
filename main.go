package main

import (
	"github.com/jnorman-us/mcfly/cmd"
	"go.minekube.com/gate/pkg/gate/config"
)

var default_config = config.DefaultConfig

func main() {
	cmd.Execute()
	// ctx := context.Background()
	// zapLogger := zap.Must(zap.NewDevelopment())
	// ctx = logr.NewContext(ctx, zapr.NewLogger(zapLogger))

	// cfg := env.Config{}
	// cfg.FlyToken = os.Getenv(env.KeyFlyToken)
	// cfg.FlyApp = "mcfly"

	// client := fly.NewFlyClient(cfg)
	// manager := mcserver.NewCloudServerManager(client)

	// manager.Initialize(ctx)

	// proxy.Plugins = append(proxy.Plugins, proxy.Plugin{
	// 	Name: "MCFlyProxy",
	// 	Init: func(ctx context.Context, proxy *proxy.Proxy) error {
	// 		return mcproxy.NewMCProxy(proxy, manager).Init(ctx)
	// 	},
	// })

	// default_config.Editions.Java.Config.Forwarding.Mode = jconfig.NoneForwardingMode
	// default_config.Editions.Java.Config.OnlineMode = false
	// gate.Start(ctx, gate.WithConfig(default_config))
}
