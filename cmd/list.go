package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"text/template"

	"github.com/fmenezes/docker-set/selector"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all environments",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		list, err := selector.List()
		if err != nil {
			log.Fatal(err)
		}
		tmpl, err := template.New("main").Parse("{{if .Active}}*{{end}}\t{{.Name}}\t{{.Driver}}\t{{if not .State}}Unknown{{else}}{{.State}}{{end}}\t{{.URL}}\n")
		if err != nil {
			log.Fatal(err)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(w, "ACTIVE\tNAME\tDRIVER\tSTATE\tURL")
		for _, entry := range list {
			tmpl.Execute(w, entry)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
