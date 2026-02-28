// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package contact

import (
	"github.com/spf13/cobra"
)

var (
	orgId     int
	token     string
	userId    string
	userEmail string
)

// Cmd represents the alias command
var Cmd = &cobra.Command{
	Use:   "contact",
	Short: "Manage user contacts",
	Long: `Use this command to manage user contacts (e-mail, phone, skype, etc.).
Cannot be executed directly, please use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("user alias called")
	//	},
}

func init() {
	Cmd.AddCommand(addCmd)
	Cmd.AddCommand(resetCmd)
	Cmd.AddCommand(lsCmd)
}
