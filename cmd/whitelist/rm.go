/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package whitelist

import (
	"fmt"
	"log"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove the list of the allowed IPs and/or CIDRs",
	Long: `Use this command to remove the list of the allowed IP addresses and/or CIDRs.
"ya360_admin:mail_write_antispam_settings" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("whitelist rm called")

		if token == "" {
			t, err := helper.GetToken()
			if err != nil {
				log.Fatalln("Incorrect settings:", err)
			}
			token = t
		}

		if orgId == 0 {
			t, err := helper.GetOrgId()
			if err != nil {
				log.Fatalln("Incorrect settings:", err)
			}
			orgId = t
		}

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mail/antispam/allowlist/ips", helper.BaseUrl, orgId)

		resp, err := helper.MakeRequest(url, "DELETE", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		log.Println("Whitelist was successfully cleared")
	},
}

func init() {
	rmCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	rmCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
}
