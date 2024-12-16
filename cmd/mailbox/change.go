/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
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

var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "Changes an existing shared mailbox",
	Long: `Use this command to change an existing shared mailbox.
"ya360_admin:mail_write_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("mailbox change called")

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
		var payload = []byte(fmt.Sprintf(`{"name":"%s", "description":"%s"}`, name, description))

		resp, err := helper.MakeRequest(url, "PUT", token, payload)
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
	},
}

func init() {
	changeCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	changeCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	changeCmd.Flags().StringVar(&name, "name", "", "shared mailbox name")
	changeCmd.Flags().StringVar(&description, "description", "", "shared mailbox description")
	changeCmd.Flags().StringVar(&mailboxId, "id", "", "shared mailbox id")

	addCmd.MarkFlagRequired("id")
}
