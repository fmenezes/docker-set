package cmd

import (
	"fmt"
	"log"

	"github.com/fmenezes/docker-set/selector"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop [name]",
	Short: "Stops the environment",
	Long: `Stops the environment. Example:

docker-set stop test`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sel, err := selector.NewSelector()
		if err != nil {
			log.Fatal(err)
		}

		err = sel.Stop(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
