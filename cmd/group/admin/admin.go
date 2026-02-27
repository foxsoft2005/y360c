/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/
package admin

import (
	"github.com/spf13/cobra"
)

var (
	orgId   int
	token   string
	groupId int
	admins  []string
)

// GroupCmd represents the group command
var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Manage Y360 group admins",
	Long: `The command provides access to admins of Y360 groups.
Cannot be executed directly, please use one if the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("group called")
	//	},
}

func init() {
	AdminCmd.AddCommand(rmCmd)
	AdminCmd.AddCommand(addCmd)
}
