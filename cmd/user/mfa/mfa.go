/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/
package mfa

import (
	"github.com/spf13/cobra"
)

var (
	orgId     int
	token     string
	userId    string
	userEmail string
)

// MfaCmd represents the 2fa command
var MfaCmd = &cobra.Command{
	Use:   "mfa",
	Short: "Manage 2fa settings for the user",
	Long: `The command provides access to two-factor auth (2fa) settings for the selected user.
Cannot be executed directly, please use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("user mfa called")
	//	},
}

func init() {
	MfaCmd.AddCommand(statusCmd)
	MfaCmd.AddCommand(resetCmd)
}
