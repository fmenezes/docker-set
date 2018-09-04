package cmd

import (
	"fmt"
	"log"

	"github.com/fmenezes/docker-set/selector"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [name]",
	Short: "Starts the environment",
	Long: `Starts the environment. Example:

docker-set start test`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sel, err := selector.NewSelector()
		if err != nil {
			log.Fatal(err)
		}

		err = sel.Start(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
