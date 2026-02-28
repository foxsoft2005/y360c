// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package disk

import (
	"github.com/spf13/cobra"
)

var (
	orgId int
	token string
)

var Cmd = &cobra.Command{
	Use:   "disk",
	Short: "Manage Y360 disk settings",
	Long: `The command provides access to Y360 disk settings.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("disk called")
	//	},
}

func init() {
	Cmd.AddCommand(logCmd)
}
