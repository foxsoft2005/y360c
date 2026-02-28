// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package domain

import (
	"github.com/spf13/cobra"
)

var (
	orgId  int
	token  string
	maxRec int
	domain string
)

// Cmd represents the domain command
var Cmd = &cobra.Command{
	Use:   "domain",
	Short: "Manage Y360 domains",
	Long: `The command provides access to Y360 domains.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("domain called")
	//	},
}

func init() {
	Cmd.AddCommand(lsCmd)
	Cmd.AddCommand(infoCmd)
}
