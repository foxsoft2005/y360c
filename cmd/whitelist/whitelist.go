// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package whitelist

import (
	"github.com/spf13/cobra"
)

var (
	orgId   int
	token   string
	allowed []string
)

var Cmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Manage Y360 anti-spam settings",
	Long: `The command provides access to Y360 anti-spam settings.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("Sorry, still under development")
	//	},
}

func init() {
	Cmd.AddCommand(lsCmd)
	Cmd.AddCommand(addCmd)
	Cmd.AddCommand(rmCmd)
}
