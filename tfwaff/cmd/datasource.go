package cmd

import (
	"github.com/openscientia/terraform-provider-atlassian/tfwaff/datasource"
	"github.com/spf13/cobra"
)

var datasourceCmd = &cobra.Command{
	Use:   "datasource",
	Short: "Generate all necessary files for a data source",
	RunE: func(cmd *cobra.Command, args []string) error {
		return datasource.Create(provider, name, force, dry_run)
	},
}

func init() {
	datasourceCmd.Flags().StringVarP(&name, "name", "n", "", "Full name of the new data-source in snake case, e.g. <provider>_<service>_<name>")
}
