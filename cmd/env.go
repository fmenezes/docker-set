package cmd

import (
	"fmt"
	"log"

	"github.com/fmenezes/docker-set/selector"
	"github.com/spf13/cobra"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env [name]",
	Short: "Sets the environment variables to an entry",
	Long: `Sets the environment variables to an entry. Example:

docker-set env test`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		env, err := selector.Env(args[0])
		if err != nil {
			log.Fatal(err)
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
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
