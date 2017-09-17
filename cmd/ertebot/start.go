package main

import (
	log "github.com/Sirupsen/logrus"

	"net/http"

	"github.com/spf13/cobra"
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

	log.Fatal(http.ListenAndServe(":8000", nil))
}
