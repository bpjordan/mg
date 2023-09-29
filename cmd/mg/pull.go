package mg

import (
	"context"
	"fmt"
	"strings"

	"github.com/bpjordan/mg/pkg/git"
	"github.com/bpjordan/mg/pkg/manifest"
	"github.com/bpjordan/mg/pkg/runtime"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pull = &cobra.Command{
	Use:              "pull",
	Aliases:          []string{"p"},
	GroupID:          "repo",
	Short:            "Pull from all repositories",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		maxConcurrent := viper.GetUint("max-connections")
		verbose := viper.GetInt("verbose")

		report, err := pullInRuntime(cmd.Context(), &man, maxConcurrent, verbose)
		if err != nil {
			return err
		}

		reportLines := make([]string, 0, 4)
		if report.Updated > 0 {
			reportLines = append(reportLines, fmt.Sprintf("%s repos updated", color.GreenString("%d", report.Updated)))
		}
		if report.NoChange > 0 {
			reportLines = append(reportLines, fmt.Sprintf("%s repos unchanged", color.GreenString("%d", report.NoChange)))
		}
		if report.Failed > 0 {
			reportLines = append(reportLines, fmt.Sprintf("%s repos encountered errors", color.RedString("%d", report.Failed)))
		}
		if report.Error > 0 {
			reportLines = append(reportLines, fmt.Sprintf("%s failed to start", color.RedString("%d", report.Error)))
		}

		fmt.Println(strings.Join(reportLines, ", "))

		return nil

	},
}

func pullInRuntime(ctx context.Context, manifest *manifest.Manifest, maxConcurrent uint, verbose int) (*git.FetchReport, error) {
	rt, err := runtime.Start(ctx, uint(len(manifest.Repos())))
	if err != nil {
		return nil, err
	}
	rt.Message = "Pulling"
	defer rt.Cleanup()

	report, err := git.Pull(rt, *manifest)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func init() {
	rootCmd.AddCommand(pull)
}
