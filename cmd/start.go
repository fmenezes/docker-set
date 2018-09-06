package cmd

import (
	"fmt"

	"github.com/fmenezes/docker-set/selector/common"
	"github.com/spf13/cobra"
)

func newStartCmd(sel common.Selector) *cobra.Command {
	return &cobra.Command{
		Use:     "start name",
		Short:   "Starts the environment",
		Example: "docker-set start test",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStart(sel, args)
		},
	}
}

func runStart(sel common.Selector, args []string) error {
	err := sel.Start(args[0])
	if err != nil {
		return err
	}
	fmt.Println("Done")
	return nil
}
