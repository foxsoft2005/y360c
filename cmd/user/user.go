/*
Copyright © 2024 Kirill Chernetsky <kirill.chernetsky@linru.ru>
*/
package user

import (
	"github.com/spf13/cobra"
	"linru.ru/y360c/cmd/user/alias"
	"linru.ru/y360c/cmd/user/mail"
	"linru.ru/y360c/cmd/user/mfa"
)

var (
	orgId  int
	token  string
	userId string
)

var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "manage Y360 users",
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
	UserCmd.AddCommand(mail.MailCmd)
}
