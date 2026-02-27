/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/
package user

import (
	"github.com/foxsoft2005/y360c/cmd/user/alias"
	"github.com/foxsoft2005/y360c/cmd/user/contact"
	"github.com/foxsoft2005/y360c/cmd/user/mfa"
	"github.com/foxsoft2005/y360c/cmd/user/rule"
	"github.com/spf13/cobra"
)

var (
	orgId     int
	token     string
	userId    string
	userEmail string
	asRaw     bool
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage Y360 users",
	Long: `The command provides access to Y360 user management.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("user called")
	//	},
}

func init() {
	UserCmd.AddCommand(listCmd)
	UserCmd.AddCommand(infoCmd)
	UserCmd.AddCommand(mfa.MfaCmd)
	UserCmd.AddCommand(alias.AliasCmd)
	UserCmd.AddCommand(rule.RuleCmd)
	UserCmd.AddCommand(contact.ContactCmd)
	UserCmd.AddCommand(setCmd)
	UserCmd.AddCommand(rmCmd)
	UserCmd.AddCommand(senderinfoCmd)
}
