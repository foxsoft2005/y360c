/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package alias

import (
	"github.com/spf13/cobra"
)

var (
	orgId  int
	token  string
	userId string
	alias  string
)

// AliasCmd represents the alias command
var AliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "manage user mailbox aliases",
	Long: `Use this command to manage user mailbox aliases.
Cannot be executed directly, please use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("user alias called")
	//	},
}

func init() {
	AliasCmd.AddCommand(addCmd)
	AliasCmd.AddCommand(rmCmd)
	AliasCmd.AddCommand(lsCmd)
}
