/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/

package mailbox

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes an shared mailbox",
	Long: `Use this command to remove an existing shared mailbox.
"ya360_admin:mail_write_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("mailbox rm called")

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

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/shared/%s", helper.BaseUrl, orgId, mailboxId)

		resp, err := helper.MakeRequest(url, "DELETE", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if resp.HttpCode != 200 {
			var errorData model.ErrorResponse
			if err := json.Unmarshal(resp.Body, &errorData); err != nil {
				log.Fatalln("Unable to evaluate data:", err)
			}
			log.Fatalf("http %d: [%d] %s", resp.HttpCode, errorData.Code, errorData.Message)
		}

		log.Println("Successfully deleted")
	},
}

func init() {
	rmCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	rmCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	rmCmd.Flags().StringVar(&mailboxId, "id", "", "shared mailbox id")

	rmCmd.MarkFlagRequired("id")
}
