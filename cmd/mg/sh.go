package mg

import (
	"fmt"
	"strings"

	"github.com/bpjordan/mg/pkg/runtime"
	"github.com/bpjordan/mg/pkg/shell"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sh = &cobra.Command{
	Use: "sh BIN (ARGS...| -- FLAGS... ARGS...)",
	Short: "Run an arbitrary shell command in all repositories",
	GroupID: "cmd",
	DisableFlagsInUseLine: true,
	Args: cobra.MinimumNArgs(1),
	TraverseChildren: true,
	RunE: shellCommand,
}

func shellCommand(cmd *cobra.Command, args []string) error {

	verbose := viper.GetInt("verbose")

	rt, err := runtime.Start(
		cmd.Context(),
		uint(len(man.Repos())),
	)

	if err != nil {
		return err
	}
	defer rt.Cleanup()

	numSuccess, numFailed, numError := shell.RunParallelCmd(
		rt,
		args[0],
		args[1:],
		man,
		verbose,
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

}

func init() {
	rootCmd.AddCommand(sh)
}
