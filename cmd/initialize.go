package cmd

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/jnorman-us/mcfly/fly"
	"github.com/jnorman-us/mcfly/mcserver"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var initializeCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision the fly.io resources for the servers, if necessary",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		zapLogger := zap.Must(zap.NewDevelopment())
		ctx = logr.NewContext(ctx, zapr.NewLogger(zapLogger))

		cfg := parseConfig()

		client := fly.NewFlyClient(cfg)
		manager := mcserver.NewCloudServerManager(client, nil, nil)

		manager.Initialize(ctx)
	},
}
