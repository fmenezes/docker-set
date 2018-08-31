package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use: "docker-set",
	Long: `docker-set is a simple tool to switch between docker environments,
    virtual machines and docker for mac`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
