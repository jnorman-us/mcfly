package cmd

import (
	"fmt"
	"os"

	"github.com/jnorman-us/mcfly/env"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var environment = env.Config{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "mcfly",
	Short: "mcfly is a Minecraft proxy that provisions fly.io machines on-demand",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	rootCmd.AddCommand(initializeCmd)
	rootCmd.AddCommand(runCmd)

	rootCmd.PersistentFlags().StringP("app", "a", "mcfly", "Name of fly.io app")

	viper.MustBindEnv("flyToken", "FLY_TOKEN")
	viper.BindPFlag("flyApp", rootCmd.PersistentFlags().Lookup("app"))
}

func parseConfig() env.Config {
	return env.Config{
		FlyToken: viper.GetString("flyToken"),
		FlyApp:   viper.GetString("flyApp"),
	}
}
