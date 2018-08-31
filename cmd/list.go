package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"text/template"

	"github.com/fmenezes/docker-set/selector"
	"github.com/fmenezes/docker-set/selector/types"
	"github.com/spf13/cobra"
)

type activeEntry struct {
	types.EnvironmentEntry
	Active bool
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all environments",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		list, err := selector.List()
		if err != nil {
			log.Fatal(err)
		}
		tmpl, err := template.New("main").Parse("{{if .Active}}*{{end}}\t{{.Name}}\t{{.Driver}}\t{{if not .State}}Unknown{{else}}{{.State}}{{end}}\n")
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(w, "ACTIVE\tNAME\tDRIVER\tSTATE")

		selected := ""
		if selector.Selected() != nil {
			selected = *selector.Selected()
		}
		for _, entry := range list {
			tmpl.Execute(w, activeEntry{
				EnvironmentEntry: entry,
				Active:           selected == entry.Name,
			})
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
