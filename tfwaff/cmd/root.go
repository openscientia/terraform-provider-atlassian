package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	provider = "atlassian"
	name     string
	force    bool
	dry_run  bool
)

var (
	rootCmd = &cobra.Command{
		Use: "tfwaff",
		Long: `tfwaff is a CLI application that generates the required files 
to implement new resources and data sources in the Terraform ATLASSIAN Provider.`,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "force creation, overwrite existing files")
	rootCmd.PersistentFlags().BoolVar(&dry_run, "dry-run", false, "do not create or overwrite files")

	rootCmd.AddCommand(resourceCmd)
	rootCmd.AddCommand(datasourceCmd)
}
