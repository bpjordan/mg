package mg

import (
	"github.com/bpjordan/mg/pkg/git"
	"github.com/bpjordan/mg/pkg/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var status = &cobra.Command{
	Use:                   "status",
	Aliases:               []string{"s"},
	GroupID:               "repo",
	Short:                 "View the current status of all repositories",
	DisableFlagsInUseLine: true,
	TraverseChildren:      true,
	RunE: func(cmd *cobra.Command, args []string) error {

		remote := viper.GetBool("remote")

		if remote {
			rt, err := runtime.Start(cmd.Context(), uint(len(man.Repos())))
			if err != nil {
				return err
			}
			defer rt.Cleanup()
			rt.Message = "Fetching"

			git.Fetch(rt, man)
		}

		return nil

	},
}

func init() {
	status.Flags().BoolP("remote", "r", false, "Show repos that are out of date with their upstream remote (requires fetching)")
	rootCmd.AddCommand(status)
}
