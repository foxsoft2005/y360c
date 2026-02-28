// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package auth

import (
	"github.com/foxsoft2005/y360c/cmd/auth/oauth"
	"github.com/spf13/cobra"
)

var (
	orgId int
	token string
)

// command definition
var Cmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage Y360 auth settings",
	Long: `The command provides access to Y360 auth settings.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("auth called")
	//	},
}

func init() {
	Cmd.AddCommand(logoutCmd)
	Cmd.AddCommand(setTtlCmd)
	Cmd.AddCommand(getTtlCmd)
	Cmd.AddCommand(oauth.OauthCmd)
}
