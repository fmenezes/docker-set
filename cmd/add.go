package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [name] [type] [location]",
	Short: "Adds a new environment entry",
	Long: `Adds a new environment entry for the list. For example:

docker-env add test vagrant /path/to/Vagrantfile
docker-env add remotehost remote tcp://123.45.67.89:2375
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return errors.New("requires 3 arguments")
		}

		if len(args[0]) == 0 {
			return errors.New("name is required")
		}

		if len(args[1]) == 0 {
			return errors.New("type is required")
		}

		if len(args[2]) == 0 {
			return errors.New("location is required")
		}

		switch args[1] {
		case "remote":
			// no special validation
		case "vagrant":
			info, err := os.Stat(args[2])
			if err != nil {
				return fmt.Errorf("Can not access %s", args[2])
			}

			if info.IsDir() {
				return errors.New("Directories are not supported, pass the Vagrantfile's full path")
			}
		default:
			return errors.New("Type should be either 'vagrant' or 'remote'")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called", args)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
