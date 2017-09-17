package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/arha/Ertebot/updater"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start bot",
	Run:   start,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func start(cmd *cobra.Command, args []string) {
	// logVersion()

	updater.Update()
}
