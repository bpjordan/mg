package mg

import (
	"fmt"
	"strings"

	"github.com/bpjordan/multigit/pkg/shell"
	"github.com/fatih/color"
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

		args = append([]string{"-c", "color.ui=always"}, args...)

		numSuccess, numFailed, numError := shell.RunParallelCmd(
			"git",
			args,
			*manifestInventory,
		)

		reportLine := make([]string, 0, 3)
		if numSuccess > 0 {
			reportLine = append(reportLine,
				fmt.Sprint(
					color.GreenString("%d", numSuccess),
					" jobs completed successfully",
			))
		}
		if numFailed > 0 {
			reportLine = append(reportLine,
				fmt.Sprint(
					color.RedString("%d", numFailed),
					" jobs exited with errors",
			))
		}
		if numError > 0 {
			reportLine = append(reportLine,
				fmt.Sprint(
					color.HiRedString("%d", numError),
					" jobs failed to start",
			))
		}

		fmt.Println(strings.Join(reportLine, ", "))
	},
}

func init() {
	rootCmd.AddCommand(git)
}
