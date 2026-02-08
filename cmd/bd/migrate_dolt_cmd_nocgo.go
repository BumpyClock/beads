package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var migrateDoltCmd = &cobra.Command{
	Use:   "dolt",
	Short: "Migrate from SQLite to Dolt backend (removed)",
	Long:  `Dolt backend has been removed. See github.com/BumpyClock/beads-dolt`,
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

func init() {
	migrateDoltCmd.Flags().Bool("dry-run", false, "Preview migration without making changes")
	migrateDoltCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt (for automation)")
	migrateCmd.AddCommand(migrateDoltCmd)
}
