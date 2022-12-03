package cmd

import (
	"github.com/openscientia/terraform-provider-atlassian/tfwaff/resource"
	"github.com/spf13/cobra"
)

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Generate all necessary files for a resource",
	RunE: func(cmd *cobra.Command, args []string) error {
		return resource.Create(provider, name, force, dry_run)
	},
}

func init() {
	resourceCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the new resource in pascal case (i.e. MixedMaps) as: <Service><Name>")
}
