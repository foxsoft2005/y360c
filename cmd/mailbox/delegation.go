/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/

package mailbox

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	enableDelegation  bool
	disableDelegation bool
)

var delegationCmd = &cobra.Command{
	Use:   "delegation",
	Short: "Enables or disables delegation for the mailbox",
	Long: `Use this command to enable or disable delegation for the existing mailbox.
"ya360_admin:mail_write_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("mailbox delegation called")

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

		var (
			url     string
			method  string
			payload []byte
		)

		if enableDelegation {
			url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/delegated", helper.BaseUrl, orgId)
			payload = []byte(fmt.Sprintf(`{"resourceId":"%s"}`, mailboxId))
			method = "PUT"
		} else {
			url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/delegated/%s", helper.BaseUrl, orgId, mailboxId)
			payload = nil
			method = "DELETE"
		}

		resp, err := helper.MakeRequest(url, method, token, payload)
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

		if enableDelegation {
			var data model.Resource
			if err := json.Unmarshal(resp.Body, &data); err != nil {
				log.Fatalln("Unable to evaluate data:", err)
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendRow(table.Row{"Resource Id", data.ResourceId})
			t.AppendRow(table.Row{"Status", "OK"})
			t.AppendSeparator()
			t.Style().Options.SeparateRows = true
			t.Render()
		} else {
			log.Println("Successfully disabled")
		}
	},
}

func init() {
	delegationCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	delegationCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	delegationCmd.Flags().StringVar(&mailboxId, "id", "", "shared mailbox id")
	delegationCmd.Flags().BoolVar(&enableDelegation, "enable", false, "enable delegation")
	delegationCmd.Flags().BoolVar(&disableDelegation, "disable", false, "disable delegation")

	delegationCmd.MarkFlagRequired("id")
	delegationCmd.MarkFlagsMutuallyExclusive("enable", "disable")
}
