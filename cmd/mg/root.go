package mg

import (
	"context"
	"fmt"

	"github.com/bpjordan/multigit/pkg/manifest"
	"github.com/spf13/cobra"
)

var (
)

var manifestContextKey struct{}

var rootCmd = &cobra.Command{
	Short: "multigit - tool for managing massive projects of multiple git repositories",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		manifestPath, err := cmd.Flags().GetString("manifest")
		if err != nil {
			return err
		}

		manifestInventory, err := manifest.ReadManifest(manifestPath)
		if err != nil {
			return err
		}

		if manifestInventory == nil {
			return fmt.Errorf("Manifest returned nil")
		}

		cmd.SetContext(context.WithValue(cmd.Context(), manifestContextKey, manifestInventory))

		return nil
	},
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	flags := rootCmd.PersistentFlags()
	flags.StringP("manifest", "m", ".mg.yml", "Path to the manifest YAML file")
	flags.UintP("max-connections", "c", 0, "Limit the number of remote operations happening concurrently")
	flags.CountP("verbose", "v", "")

	rootCmd.AddGroup(&cobra.Group{ID: "repo", Title: "Manage Repositories"})
	rootCmd.AddGroup(&cobra.Group{ID: "cmd", Title: "Run Arbitrary Commands"})
}
