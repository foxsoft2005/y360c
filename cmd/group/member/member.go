/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/
package member

import (
	"errors"

	"github.com/spf13/cobra"
)

type groupMemberType string

const userMember groupMemberType = "user"

func (e *groupMemberType) String() string {
	return string(*e)
}

func (e *groupMemberType) Set(v string) error {
	switch v {
	case "user", "group", "department":
		*e = groupMemberType(v)
		return nil
	default:
		return errors.New(`must be one of "user", "group" or "department"`)
	}
}

func (e *groupMemberType) Type() string {
	return "groupMemberType"
}

var (
	orgId      int
	token      string
	groupId    int
	memberId   string
	memberType groupMemberType
)

// GroupCmd represents the group command
var MemberCmd = &cobra.Command{
	Use:   "member",
	Short: "Manage Y360 user groups",
	Long: `The command provides access to Y360 groups.
Cannot be executed directly, please use one if the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("group called")
	//	},
}

func init() {
	MemberCmd.AddCommand(lsCmd)
	MemberCmd.AddCommand(rmCmd)
	MemberCmd.AddCommand(addCmd)
}
