// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package rule

import (
	"github.com/spf13/cobra"
)

var (
	orgId  int
	token  string
	userId string
	ruleId string
)

// Cmd represents the rule command
var Cmd = &cobra.Command{
	Use:   "rule",
	Short: "Manage mailbox rules for the user",
	Long: `The command provides access to mailbox rules (forwards and autoreplies) for the selected user.
Cannot be executed directly, please use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("user mail called")
	//	},
}

func init() {
	Cmd.AddCommand(lsCmd)
}
