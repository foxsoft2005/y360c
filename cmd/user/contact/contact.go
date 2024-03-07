/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package contact

import (
	"github.com/spf13/cobra"
)

var (
	orgId  int
	token  string
	userId string
)

// AliasCmd represents the alias command
var ContactCmd = &cobra.Command{
	Use:   "contact",
	Short: "manage user contacts",
	Long: `Use this command to manage user contacts (e-mail, phone, skype, etc.).
Cannot be executed directly, please use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("user alias called")
	//	},
}

func init() {
	ContactCmd.AddCommand(addCmd)
	ContactCmd.AddCommand(resetCmd)
	ContactCmd.AddCommand(lsCmd)
}
