package cmd

import (
	"fmt"

	"github.com/fmenezes/docker-set/selector/common"
	"github.com/spf13/cobra"
)

func newAddCmd(sel common.Selector) *cobra.Command {
	return &cobra.Command{
		Use:     "add name driver location",
		Short:   "Adds a new environment entry",
		Example: "docker-env add test vagrant /path/to/Vagrantfile",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(sel, args)
		},
	}
}

func runAdd(sel common.Selector, args []string) error {
	err := sel.Add(common.EnvironmentEntry{
		Name:     args[0],
		Driver:   args[1],
		Location: &args[2],
	})
	if err != nil {
		return err
	}
	fmt.Println("Done")
	return nil
}
