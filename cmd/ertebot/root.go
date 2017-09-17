package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ertebot <subcommand>",
	Short: "Telegram bot",
	Run:   nil,
}

func init() {
	cobra.OnInitialize()
}
