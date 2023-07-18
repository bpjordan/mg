package mg

import (
	"fmt"
	"strings"
	"time"

	"github.com/bpjordan/multigit/pkg/manifest"
	"github.com/bpjordan/multigit/pkg/runtime"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var debug = &cobra.Command{
	Use: "debug",
	Hidden: true,
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Debug args: ", strings.Join(args, " "))
		for _, arg := range args {
			switch arg {
			case "manifest":
				manifest := cmd.Context().Value(manifestContextKey).(*manifest.Manifest)
				spew.Dump(*manifest)
			case "statusline":
				maxConcurrent, err := cmd.Flags().GetUint("max-connections")

				rt, err := runtime.Start(cmd.Context(), 1, maxConcurrent)
				if err != nil {
					fmt.Fprintln(cmd.OutOrStderr(), err)
				}
				defer rt.Cleanup()

				for i := 1; i <= 50; i++ {
					select {
					case <-rt.Finished():
						break
					default:
						fmt.Println(i)
						time.Sleep(time.Second)
					}
				}
			default:
				fmt.Printf("Unknown debug parameter `%s`\n", arg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(debug)
}
