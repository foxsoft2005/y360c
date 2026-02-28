// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package oauth

import (
	"github.com/spf13/cobra"
)

var (
	orgId int
	token string
)

// command definition
var OauthCmd = &cobra.Command{
	Use:   "oauth",
	Short: "Manage Y360 OAuth settings",
	Long: `The command provides access to Y360 OAuth settings.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("auth called")
	//	},
}

func init() {
	OauthCmd.AddCommand(statusCmd)
}
