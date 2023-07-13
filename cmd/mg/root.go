package mg

import (
	"fmt"
	"os"

	"github.com/bpjordan/multigit/pkg/manifest"
	"github.com/spf13/cobra"
)

var (
	manifestPath string
	manifestInventory *manifest.Manifest
	maxConcurrent uint
)

var rootCmd = &cobra.Command{
	Short: "multigit - tool for managing massive projects of multiple git repositories",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		var err error
		manifestInventory, err = manifest.ReadManifest(manifestPath)
		if err != nil {
			return err
		}

		if manifestInventory == nil {
			return fmt.Errorf("Manifest returned nil")
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&manifestPath, "manifest", "m", ".mg.yml", "Path to the manifest YAML file")
	rootCmd.PersistentFlags().UintVarP(&maxConcurrent, "max-connections", "c", 0, "Limit the number of remote operations happening concurrently")

	rootCmd.AddGroup(&cobra.Group{ID: "cmd", Title: "Run Arbitrary Commands"})
}
