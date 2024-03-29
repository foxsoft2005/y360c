/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
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
	Short: "manage Y360 auth settings",
	Long: `The command provides access to Y360 auth settings.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("auth called")
	//	},
}

func init() {
	AuthCmd.AddCommand(logoutCmd)
	AuthCmd.AddCommand(setCookieTtlCmd)
	AuthCmd.AddCommand(getCookieTtlCmd)
}
