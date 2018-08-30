package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "Removes an environment entry",
	Long: `Removes an environment entry from the list. Example:

docker-set rm test`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rm called", args)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
