package mg

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/davecgh/go-spew/spew"
)

var debug = &cobra.Command{
	Use: "debug",
	Hidden: true,
	DisableFlagParsing: true,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Debug args: ", strings.Join(args, " "))
		spew.Dump(manifestInventory)
	},
}

func init() {
	rootCmd.AddCommand(debug)
}
