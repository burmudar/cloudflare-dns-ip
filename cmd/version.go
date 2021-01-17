package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var BuildTime, BuildVersion, BuildCommit string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "displays version information",
	Long:  "displays the Build Time and Version of this binary",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Build Time: %s\nVersion: %s\nCommit: %s\n", BuildTime, BuildVersion, BuildCommit)
	},
}
