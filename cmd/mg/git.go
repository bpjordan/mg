package mg

import (
	"fmt"
	"strings"

	"github.com/bpjordan/multigit/pkg/shell"
	"github.com/spf13/cobra"
)

var git = &cobra.Command{
	Use: "git ARG ...",
	Aliases: []string{"g"},
	Short: "Run an arbitrary `git` command in all repositories",
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Git args: ", strings.Join(args, " "))
		shell.RunCmd("git", args, manifestInventory.Paths())
	},
}

func init() {
	rootCmd.AddCommand(git)
}
