package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print mdview version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mdview v%s\n", Version)
	},
}
