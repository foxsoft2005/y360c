/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package rule

import (
	"fmt"
	"log"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes mailbox rule by id",
	Long: `Use this command to remove mailbox rule (autoreply or forward) for the selected user.
"ya360_admin:mail_write_user_settings" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user mail rules called")

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

		if !helper.Confirm("Do you REALLY want to DELETE the selected mailbox rule (y[es]|no)?") {
			log.Fatalln("Aborted by the user")
		}

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mail/users/%s/settings/user_rules/%s", helper.BaseUrl, orgId, userId, ruleId)

		resp, err := helper.MakeRequest(url, "DELETE", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		log.Printf("Rule %s successfully deleted", ruleId)
	},
}

func init() {
	rmCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	rmCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	rmCmd.Flags().StringVar(&userId, "id", "", "user id")
	rmCmd.Flags().StringVar(&ruleId, "rule-id", "", "rule id")

	rmCmd.MarkFlagRequired("id")
	rmCmd.MarkFlagRequired("rule-id")
}
