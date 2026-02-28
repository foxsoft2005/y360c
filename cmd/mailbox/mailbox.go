// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

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
	mailboxName string
)

var Cmd = &cobra.Command{
	Use:   "mailbox",
	Short: "Manage Y360 shared and delegated mailboxes",
	Long: `The command provides access to Y360 shared and delegated mailboxes.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("Sorry, still under development")
	//	},
}

func init() {
	Cmd.AddCommand(lsCmd)
	Cmd.AddCommand(infoCmd)
	Cmd.AddCommand(addCmd)
	Cmd.AddCommand(rmCmd)
	Cmd.AddCommand(delegationCmd)
	Cmd.AddCommand(hasAccessCmd)
	Cmd.AddCommand(accessCmd)
	Cmd.AddCommand(sharedWithCmd)
	Cmd.AddCommand(statusCmd)
	Cmd.AddCommand(changeCmd)
}
