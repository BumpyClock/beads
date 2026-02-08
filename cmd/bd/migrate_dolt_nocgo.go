package main

import (
	"fmt"
	"os"
)

// handleToDoltMigration is a stub; Dolt backend has been removed.
func handleToDoltMigration(dryRun bool, autoYes bool) {
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

// handleToSQLiteMigration is a stub; Dolt backend has been removed.
func handleToSQLiteMigration(dryRun bool, autoYes bool) {
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


