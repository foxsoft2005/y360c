/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package org

import (
	"github.com/spf13/cobra"
)

// command level flags
var (
	token  string
	maxRec int
)

// orgCmd represents the org command
var OrgCmd = &cobra.Command{
	Use:   "org",
	Short: "Manage Y360 organizations",
	Long: `The command provides access to Y360 organizations.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("org called")
	//	},
}

func init() {
	OrgCmd.AddCommand(listCmd)
}
