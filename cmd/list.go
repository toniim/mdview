package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/bilabl/mdview/internal/process"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all running mdview servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		instances := process.FindRunningInstances()
		if len(instances) == 0 {
			fmt.Println("No running instances")
			return nil
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 4, 2, ' ', 0)
		fmt.Fprintln(w, "PID\tHOST\tPORT\tPATH")
		for _, inst := range instances {
			host := inst.Host
			if host == "" {
				host = "-"
			}
			path := inst.Path
			if path == "" {
				path = "-"
			}
			fmt.Fprintf(w, "%d\t%s\t%d\t%s\n", inst.Pid, host, inst.Port, path)
		}
		w.Flush()
		return nil
	},
}
