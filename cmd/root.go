// cmd package contains cli logic, reads command line arguments and respond to STDOUT
package cmd

import (
	"log"
	"os"

	"github.com/fmenezes/docker-set/selector"
	"github.com/fmenezes/docker-set/selector/common"
	"github.com/fmenezes/docker-set/selector/drivers/dockerformac"
	"github.com/fmenezes/docker-set/selector/drivers/dockermachine"
	"github.com/fmenezes/docker-set/selector/drivers/vagrant"
	"github.com/fmenezes/docker-set/selector/storage"
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use: "docker-set",
		Long: `docker-set is a simple tool to switch between docker environments,
virtual machines and docker for mac`,
	}
}

func buildSelector() (common.Selector, error) {
	sel := selector.NewSelector()

	dockerForMacDriver := dockerformac.NewDriver()
	if dockerForMacDriver.IsSupported() {
		sel.RegisterDriver(dockerForMacDriver)
	}

	dockerMachineDriver := dockermachine.NewDriver()
	if dockerMachineDriver.IsSupported() {
		sel.RegisterDriver(dockerMachineDriver)
	}

	file, err := storage.GetFilePath()
	if err != nil {
		return nil, err
	}
	vagrantDriver := vagrant.NewDriver(storage.NewFileStorage(file))
	if vagrantDriver.IsSupported() {
		sel.RegisterDriver(vagrantDriver)
	}

	return sel, nil
}

// Executes the main cli logic
func Execute() {
	sel, err := buildSelector()
	if err != nil {
		log.Fatal(err)
	}

	rootCmd := newRootCmd()
	rootCmd.AddCommand(newAddCmd(sel))
	rootCmd.AddCommand(newRemoveCmd(sel))
	rootCmd.AddCommand(newListCmd(sel))
	rootCmd.AddCommand(newEnvCmd(sel))
	rootCmd.AddCommand(newStartCmd(sel))
	rootCmd.AddCommand(newStopCmd(sel))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
