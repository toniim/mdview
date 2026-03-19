package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mdview",
	Short: "A calm, book-like viewer for markdown files",
	Long:  "mdview renders markdown files with a warm, novel-reader UI.\nDrop it anywhere and run — single binary, no dependencies.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(versionCmd)
}
