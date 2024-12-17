/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package mfa

import (
	"github.com/spf13/cobra"
)

var (
	orgId int
	token string
)

// mfaCmd represents the mfa command
var MfaCmd = &cobra.Command{
	Use:   "mfa",
	Short: "Manage Y360 2FA settings",
	Long: `The command provides access to Y360 two-factor auth (2FA) settings.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("mfa called")
	//	},
}

func init() {
	MfaCmd.AddCommand(statusCmd)
	MfaCmd.AddCommand(enableCmd)
	MfaCmd.AddCommand(disableCmd)
}
