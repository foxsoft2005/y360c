/*
Copyright © 2024 Kirill Chernetsky <kirill.chernetsky@linru.ru>
*/
package dns

import (
	"github.com/spf13/cobra"
)

var (
	orgId  int
	token  string
	maxRec int
	domain string
)

// DnsCmd represents the dns command
var DnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "manage DNS records of the domain",
	Long: `The command provides access to DNS records of the domain.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("dns called")
	//	},
}

func init() {
	DnsCmd.AddCommand(lsCmd)
}
