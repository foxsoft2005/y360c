// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package dept

import (
	"github.com/foxsoft2005/y360c/cmd/dept/alias"
	"github.com/foxsoft2005/y360c/cmd/dept/member"
	"github.com/spf13/cobra"
)

var (
	orgId       int
	token       string
	name        string
	description string
	label       string
	deptId      int
	parentId    int
	externalId  string
	headId      string
)

// deptCmd represents the dept command
var DeptCmd = &cobra.Command{
	Use:   "dept",
	Short: "Manage Y360 departments",
	Long: `The command provides access to Y360 departments.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("dept called")
	//	},
}

func init() {
	DeptCmd.AddCommand(lsCmd)
	DeptCmd.AddCommand(infoCmd)
	DeptCmd.AddCommand(addCmd)
	DeptCmd.AddCommand(rmCmd)
	DeptCmd.AddCommand(changeCmd)
	DeptCmd.AddCommand(alias.Cmd)
	DeptCmd.AddCommand(member.Cmd)
}
