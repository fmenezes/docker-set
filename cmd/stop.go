package cmd

import (
	"fmt"

	"github.com/fmenezes/docker-set/selector/common"
	"github.com/spf13/cobra"
)

func newStopCmd(sel common.Selector) *cobra.Command {
	return &cobra.Command{
		Use:     "stop name",
		Short:   "Stops the environment",
		Example: "docker-set stop test",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStop(sel, args)
		},
	}
}

func runStop(sel common.Selector, args []string) error {
	err := sel.Stop(args[0])
	if err != nil {
		return err
	}
	fmt.Println("Done")
	return nil
}
