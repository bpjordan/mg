package mg

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var debug = &cobra.Command{
	Use: "debug",
	Hidden: true,
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Debug args: ", strings.Join(args, " "))
		if manifestInventory != nil {
			fmt.Println("Manifest: ", *manifestInventory)
		}
	},
}

func init() {
	rootCmd.AddCommand(debug)
}
