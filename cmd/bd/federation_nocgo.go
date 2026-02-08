package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var federationCmd = &cobra.Command{
	Use:     "federation",
	GroupID: "sync",
	Short:   "Manage peer-to-peer federation (removed)",
	Long: `Dolt backend has been removed. See github.com/BumpyClock/beads-dolt

Federation required the Dolt storage backend, which is no longer included
in this binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Dolt backend has been removed.")
		fmt.Println("See github.com/BumpyClock/beads-dolt")
	},
}

func init() {
	rootCmd.AddCommand(federationCmd)
}


