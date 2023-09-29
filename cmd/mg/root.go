package mg

import (
	"github.com/bpjordan/mg/pkg/manifest"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	man manifest.Manifest
)

var rootCmd = &cobra.Command{
	Short: "multigit - tool for managing massive projects of multiple git repositories",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

		viper.BindPFlags(cmd.Flags())
		viper.SetConfigName(".mg")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.config/mg")
		viper.AddConfigPath(".")

		if manifestPath, _ := cmd.Flags().GetString("manifest"); manifestPath != "" {
			viper.SetConfigFile(manifestPath)
		}

		err := viper.ReadInConfig()
		if err != nil {
			return err
		}

		err = viper.Unmarshal(&man, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc()))
		if err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	flags := rootCmd.PersistentFlags()
	flags.StringP("manifest", "m", "", "Path to the manifest YAML file")
	flags.UintP("max-connections", "c", 0, "Limit the number of simultaneous connections to the upstream repository (useful if your remote utilizes rate limiting)")
	flags.CountP("verbose", "v", "")

	rootCmd.AddGroup(&cobra.Group{ID: "repo", Title: "Manage Repositories"})
	rootCmd.AddGroup(&cobra.Group{ID: "cmd", Title: "Run Arbitrary Commands"})
}
