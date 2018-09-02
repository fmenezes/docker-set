package cmd

import (
	"fmt"
	"log"

	"github.com/fmenezes/docker-set/selector"
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
		sel, err := selector.NewSelector()
		if err != nil {
			log.Fatal(err)
		}

		err = sel.Remove(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
