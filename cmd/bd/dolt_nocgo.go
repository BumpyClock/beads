package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func noCGODoltError(cmd *cobra.Command, args []string) {
	if jsonOutput {
		outputJSON(map[string]interface{}{
			"error":   "dolt_not_available",
			"message": "Dolt backend has been removed. See github.com/BumpyClock/beads-dolt",
		})
	} else {
		fmt.Fprintf(os.Stderr, "Error: Dolt backend has been removed.\n")
		fmt.Fprintf(os.Stderr, "See github.com/BumpyClock/beads-dolt\n")
	}
	os.Exit(1)
}

var doltCmd = &cobra.Command{
	Use:     "dolt",
	GroupID: "setup",
	Short:   "Configure Dolt database settings",
	Long:    `Dolt backend has been removed. See github.com/BumpyClock/beads-dolt`,
	Run:     noCGODoltError,
}

var doltShowCmdNoCGO = &cobra.Command{
	Use:   "show",
	Short: "Show current Dolt configuration with connection status",
	Run:   noCGODoltError,
}

var doltSetCmdNoCGO = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a Dolt configuration value",
	Args:  cobra.ExactArgs(2),
	Run:   noCGODoltError,
}

var doltTestCmdNoCGO = &cobra.Command{
	Use:   "test",
	Short: "Test connection to Dolt server",
	Run:   noCGODoltError,
}

func init() {
	doltSetCmdNoCGO.Flags().Bool("update-config", false, "Also write to config.yaml for team-wide defaults")
	doltCmd.AddCommand(doltShowCmdNoCGO)
	doltCmd.AddCommand(doltSetCmdNoCGO)
	doltCmd.AddCommand(doltTestCmdNoCGO)
	rootCmd.AddCommand(doltCmd)
}
