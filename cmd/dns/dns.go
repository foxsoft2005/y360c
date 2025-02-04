/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
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
	Short: "Manage Y360 DNS records of the domain",
	Long: `The command provides access to DNS records of the domain.
Cannot be executed directly, use one of the available sub-commands.`,
	//	Run: func(cmd *cobra.Command, args []string) {
	//		fmt.Println("dns called")
	//	},
}

func init() {
	DnsCmd.AddCommand(lsCmd)
}
