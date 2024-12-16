/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/

package mailbox

import (
	"github.com/spf13/cobra"
)

var (
	orgId       int
	token       string
	mailboxId   string
	name        string
	description string
	userId      string
)

var MailboxCmd = &cobra.Command{
	Use:   "mailbox",
	Short: "Manage Y360 shared and delegated mailboxes",
	Long: `The command provides access to Y360 shared and delegated mailboxes.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("Sorry, still under development")
	//	},
}

func init() {
	MailboxCmd.AddCommand(lsCmd)
	MailboxCmd.AddCommand(infoCmd)
	MailboxCmd.AddCommand(addCmd)
	MailboxCmd.AddCommand(rmCmd)
	MailboxCmd.AddCommand(delegationCmd)
	MailboxCmd.AddCommand(hasAccessCmd)
	MailboxCmd.AddCommand(setaccessCmd)
	MailboxCmd.AddCommand(sharedWithCmd)
	MailboxCmd.AddCommand(statusCmd)
	MailboxCmd.AddCommand(changeCmd)
}
