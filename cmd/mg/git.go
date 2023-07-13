package mg

import (
	"fmt"
	"strings"

	"github.com/bpjordan/multigit/pkg/runtime"
	"github.com/bpjordan/multigit/pkg/shell"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var git = &cobra.Command{
	Use: "git {ARGS...| -- FLAGS... ARGS...}",
	Aliases: []string{"g"},
	Short: "Run an arbitrary `git` command in all repositories",
	DisableFlagsInUseLine: true,
	Args: cobra.MinimumNArgs(1),
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		args = append([]string{"-c", "color.ui=always"}, args...)

		rt, err := runtime.Start(
			cmd.Context(),
			uint(len(manifestInventory.Repos)),
			maxConcurrent,
		)
		if err != nil {
			return err
		}
		defer rt.Cleanup()

		numSuccess, numFailed, numError := shell.RunParallelCmd(
			rt,
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

		return nil
	},
}

func init() {
	rootCmd.AddCommand(git)
}
