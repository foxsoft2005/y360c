/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package group

import (
	"github.com/foxsoft2005/y360c/cmd/group/admin"
	"github.com/foxsoft2005/y360c/cmd/group/member"
	"github.com/spf13/cobra"
)

var (
	orgId       int
	token       string
	groupId     int
	name        string
	label       string
	description string
	externalId  string
	admins      []string
)

// GroupCmd represents the group command
var GroupCmd = &cobra.Command{
	Use:   "group",
	Short: "manage Y360 user groups",
	Long: `The command provides access to Y360 groups.
Cannot be executed directly, please use one if the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("group called")
	//	},
}

func init() {
	GroupCmd.AddCommand(infoCmd)
	GroupCmd.AddCommand(lsCmd)
	GroupCmd.AddCommand(addCmd)
	GroupCmd.AddCommand(rmCmd)
	GroupCmd.AddCommand(member.MemberCmd)
	GroupCmd.AddCommand(admin.AdminCmd)
}
