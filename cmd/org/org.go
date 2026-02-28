// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package org

import (
	"github.com/spf13/cobra"
)

// command level flags
var (
	token  string
	maxRec int
)

// Cmd represents the org command
var Cmd = &cobra.Command{
	Use:   "org",
	Short: "Manage Y360 organizations",
	Long: `The command provides access to Y360 organizations.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("org called")
	//	},
}

func init() {
	Cmd.AddCommand(listCmd)
}
