// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package alias

import (
	"github.com/spf13/cobra"
)

var (
	orgId  int
	token  string
	deptId int
	alias  string
)

// AliasCmd represents the alias command
var Cmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage user mailbox aliases",
	Long: `Use this command to manage user mailbox aliases.
Cannot be executed directly, please use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("user alias called")
	//	},
}

func init() {
	Cmd.AddCommand(addCmd)
	Cmd.AddCommand(rmCmd)
	Cmd.AddCommand(lsCmd)
}
