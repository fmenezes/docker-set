package cmd

import (
	"fmt"

	"github.com/fmenezes/docker-set/selector/common"
	"github.com/spf13/cobra"
)

func newEnvCmd(sel common.Selector) *cobra.Command {
	return &cobra.Command{
		Use:     "env name",
		Short:   "Sets the environment variables to an entry",
		Example: "docker-set env test",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEnv(sel, args)
		},
	}
}

func runEnv(sel common.Selector, args []string) error {
	env, err := sel.Env(args[0])
	if err != nil {
		return err
	}
	for key, value := range env {
		if value == nil {
			fmt.Printf("export %s=\n", key)
		} else {
			fmt.Printf("export %s='%s'\n", key, *value)
		}
	}
	fmt.Println("# Run this command to configure your shell:")
	fmt.Printf("# eval $(docker-set env %s)\n", args[0])
	return nil
}
