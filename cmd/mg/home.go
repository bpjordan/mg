package mg

import (
	"github.com/bpjordan/multigit/pkg/git"
	"github.com/bpjordan/multigit/pkg/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var home = &cobra.Command{
	Use:                   "home",
	GroupID:               "repo",
	Short:                 "View the current status of all repositories",
	DisableFlagsInUseLine: true,
	TraverseChildren:      true,
	RunE: func(cmd *cobra.Command, args []string) error {

		sync := viper.GetBool("sync")

		rt, err := runtime.Start(cmd.Context(), uint(len(man.Repos())))
		if err != nil {
			return err
		}
		defer rt.Cleanup()
		rt.Message = "Checking Out Home"

		git.Checkout(nil, rt, man)

		if sync {
			pull.RunE(cmd, args)
		}

		return nil

	},
}

func init() {
	home.Flags().BoolP("sync", "s", false, "Pull latest changes after checking out home branch")
	rootCmd.AddCommand(home)
}
