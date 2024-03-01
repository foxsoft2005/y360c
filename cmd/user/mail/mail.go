/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package mail

import (
	"github.com/spf13/cobra"
)

var (
	orgId  int
	token  string
	userId string
)

// MailCmd represents the mail command
var MailCmd = &cobra.Command{
	Use:   "mail",
	Short: "manage mail settings for the user",
	Long: `The command provides access to mail settings (sender info, mail rules etc) for the selected user.
Cannot be executed directly, please use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("user mail called")
	//	},
}

func init() {
	MailCmd.AddCommand(rulesCmd)
	MailCmd.AddCommand(senderinfoCmd)
}
