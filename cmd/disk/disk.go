/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package disk

import (
	"github.com/spf13/cobra"
)

var (
	orgId int
	token string
)

var DiskCmd = &cobra.Command{
	Use:   "disk",
	Short: "Manage Y360 disk settings",
	Long: `The command provides access to Y360 disk settings.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("disk called")
	//	},
}

func init() {
	DiskCmd.AddCommand(logCmd)
}
