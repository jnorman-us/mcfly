package cmd

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/jnorman-us/mcfly/fly"
	"github.com/jnorman-us/mcfly/mcproxy"
	"github.com/jnorman-us/mcfly/mcserver"
	"github.com/spf13/cobra"
	jconfig "go.minekube.com/gate/pkg/edition/java/config"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"go.minekube.com/gate/pkg/gate"
	"go.minekube.com/gate/pkg/gate/config"
	"go.uber.org/zap"
)

var default_config = config.DefaultConfig

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the proxy",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		zapLogger := zap.Must(zap.NewDevelopment())
		ctx = logr.NewContext(ctx, zapr.NewLogger(zapLogger))

		cfg := parseConfig()

		client := fly.NewFlyClient(cfg)
		manager := mcserver.NewCloudServerManager(client)

		err := manager.FindCloudResources(ctx)
		if err != nil {
			zapLogger.Panic("Problem finding existing cloud resources", zap.Error(err))
		}

		proxy.Plugins = append(proxy.Plugins, proxy.Plugin{
			Name: "MCFlyProxy",
			Init: func(ctx context.Context, proxy *proxy.Proxy) error {
				return mcproxy.NewMCProxy(proxy, manager).Init(ctx)
			},
		})

		default_config.Editions.Java.Config.Forwarding.Mode = jconfig.NoneForwardingMode
		default_config.Editions.Java.Config.OnlineMode = false
		gate.Start(ctx, gate.WithConfig(default_config))
	},
}
