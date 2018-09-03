// cmd package contains cli logic, reads command line arguments and respond to STDOUT
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

// Executes the main cli logic
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
