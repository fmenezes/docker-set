package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"

	"github.com/fmenezes/docker-set/selector/common"
	"github.com/spf13/cobra"
)

type activeEntry struct {
	common.EnvironmentEntryWithState
	Active bool
}

func newListCmd(sel common.Selector) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists all environments",
		Example: "docker-set list",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(sel, args)
		},
	}
}

func runList(sel common.Selector, args []string) error {
	tmpl, err := template.New("main").Parse("{{if .Active}}*{{end}}\t{{.Name}}\t{{.Driver}}\t{{if not .State}}Unknown{{else}}{{.State}}{{end}}\n")
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "ACTIVE\tNAME\tDRIVER\tSTATE")

	selected := ""
	if sel.Selected() != nil {
		selected = *sel.Selected()
	}
	for entry := range sel.List() {
		tmpl.Execute(w, activeEntry{
			EnvironmentEntryWithState: entry,
			Active: selected == entry.Name,
		})
	}
	w.Flush()
	return nil
}
