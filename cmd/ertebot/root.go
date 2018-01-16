package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/arha/Ertebot/configuration"
)

var rootCmd = &cobra.Command{
	Use:   "ertebot <subcommand>",
	Short: "Telegram bot",
	Run:   nil,
}

func init() {
	cobra.OnInitialize(func() {
		configuration.LoadConfig()
	})
}
