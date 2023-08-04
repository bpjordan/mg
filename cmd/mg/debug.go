package mg

import (
	"fmt"
	"strings"
	"time"

	"github.com/bpjordan/multigit/pkg/runtime"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debug = &cobra.Command{
	Use: "debug",
	Hidden: true,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Debug args: ", strings.Join(args, " "))
		for _, arg := range args {
			switch arg {
			case "manifest":
				spew.Dump(man)
			case "statusline":
				rt, err := runtime.Start(cmd.Context(), 1)
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
			case "config":
				spew.Dump(viper.AllSettings())
			default:
				fmt.Printf("Unknown debug parameter `%s`\n", arg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(debug)
}
