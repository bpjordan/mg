package mg

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var git = &cobra.Command{
	Use: "git {ARGS...| -- FLAGS... ARGS...}",
	Aliases: []string{"g"},
	GroupID: "cmd",
	Short: "Run an arbitrary `git` command in all repositories",
	DisableFlagsInUseLine: true,
	Args: cobra.MinimumNArgs(1),
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if color.NoColor {
			args = append([]string{"git"}, args...)
		} else {
			args = append([]string{"git", "-c", "color.ui=always"}, args...)
		}

		return shellCommand(cmd, args)

	},
}

func init() {
	rootCmd.AddCommand(git)
}
