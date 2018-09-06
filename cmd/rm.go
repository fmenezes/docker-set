package cmd

import (
	"fmt"

	"github.com/fmenezes/docker-set/selector/common"
	"github.com/spf13/cobra"
)

func newRemoveCmd(sel common.Selector) *cobra.Command {
	return &cobra.Command{
		Use:     "rm name",
		Short:   "Removes an environment entry",
		Example: "docker-set rm test",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRemove(sel, args)
		},
	}
}

func runRemove(sel common.Selector, args []string) error {
	err := sel.Remove(args[0])
	if err != nil {
		return err
	}
	fmt.Println("Done")
	return nil
}
