/*
Copyright © 2024 Kirill Chernetsky <kirill.chernetsky@linru.ru>
*/
package domain

import (
	"github.com/spf13/cobra"
)

var (
	orgId  int
	token  string
	maxRec int
	domain string
)

// DomainCmd represents the domain command
var DomainCmd = &cobra.Command{
	Use:   "domain",
	Short: "manage Y360 domains",
	Long: `The command provides access to Y360 domains.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("domain called")
	//	},
}

func init() {
	DomainCmd.AddCommand(lsCmd)
	DomainCmd.AddCommand(infoCmd)
}
