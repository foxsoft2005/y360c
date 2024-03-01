/*
Copyright © 2024 Kirill Chernetsky <kirill.chernetsky@linru.ru>
*/
package whitelist

import (
	"github.com/spf13/cobra"
)

var (
	orgId   int
	token   string
	allowed []string
)

var WhitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "manage Y360 anti-spam settings",
	Long: `The command provides access to Y360 anti-spam settings.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("Sorry, still under development")
	//	},
}

func init() {
	WhitelistCmd.AddCommand(lsCmd)
	WhitelistCmd.AddCommand(addCmd)
	WhitelistCmd.AddCommand(rmCmd)
}
