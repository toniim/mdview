package cmd

import (
	"fmt"

	"github.com/bilabl/mdview/internal/process"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all running mdview servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		instances := process.FindRunningInstances()
		if len(instances) == 0 {
			fmt.Println("No server running to stop")
			return nil
		}
		stopped := process.StopAllServers()
		fmt.Printf("Stopped %d server(s)\n", stopped)
		return nil
	},
}
