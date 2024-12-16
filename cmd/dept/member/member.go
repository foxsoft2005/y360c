/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package member

import (
	"github.com/spf13/cobra"
)

var (
	orgId  int
	token  string
	deptId int
)

// GroupCmd represents the group command
var MemberCmd = &cobra.Command{
	Use:   "member",
	Short: "Manage Y360 department members",
	Long: `The command provides access to Y360 department members.
Cannot be executed directly, please use one if the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("group called")
	//	},
}

func init() {
	MemberCmd.AddCommand(lsCmd)
}
