/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package auth

import (
	"github.com/spf13/cobra"
)

var (
	orgId int
	token string
)

// command definition
var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage Y360 auth settings",
	Long: `The command provides access to Y360 auth settings.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("auth called")
	//	},
}

func init() {
	AuthCmd.AddCommand(logoutCmd)
	AuthCmd.AddCommand(setTtlCmd)
	AuthCmd.AddCommand(getTtlCmd)
}
