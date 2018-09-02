package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/fmenezes/docker-set/selector"
	"github.com/fmenezes/docker-set/selector/common"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name] [driver] [location]",
	Short: "Adds a new environment entry",
	Long: `Adds a new environment entry for the list. For example:

docker-env add test vagrant /path/to/Vagrantfile
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return errors.New("requires 3 arguments")
		}

		if len(args[0]) == 0 {
			return errors.New("name is required")
		}

		if len(args[1]) == 0 {
			return errors.New("driver is required")
		}

		if len(args[2]) == 0 {
			return errors.New("location is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sel, err := selector.NewSelector()
		if err != nil {
			log.Fatal(err)
		}

		err = sel.Add(common.EnvironmentEntry{
			Name:     args[0],
			Driver:   args[1],
			Location: &args[2],
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
