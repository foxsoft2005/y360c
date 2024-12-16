/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package mailbox

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	describeResource bool
)

// hasAccessCmd represents the hasAccess command
var hasAccessCmd = &cobra.Command{
	Use:   "hasAccess",
	Short: "Gets mailboxes that user has access to",
	Long: `Use this command to get mailboxes that selected user has access to.
"ya360_admin:mail_read_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user mail hasAccess called")

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

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/resources/%s", helper.BaseUrl, orgId, userId)

		resp, err := helper.MakeRequest(url, "GET", token, nil)
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

		var data model.ResourceList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		if describeResource {
			t.AppendHeader(table.Row{"Id", "Email", "Type", "Roles"})
			for _, e := range data.Items {
				if e.ResourceType == "delegated" {
					user, _ := helper.GetUserById(orgId, token, e.ResourceId)
					t.AppendRow(table.Row{e.ResourceId, user.Email, e.ResourceType, strings.Join(e.Items, ",")})
				} else {
					mailbox, _ := helper.GetMailboxById(orgId, token, e.ResourceId)
					t.AppendRow(table.Row{e.ResourceId, mailbox.Email, e.ResourceType, strings.Join(e.Items, ",")})
				}
			}
		} else {
			t.AppendHeader(table.Row{"Id", "Type", "Roles"})
			for _, e := range data.Items {
				t.AppendRow(table.Row{e.ResourceId, e.ResourceType, strings.Join(e.Items, ",")})
			}
		}
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	hasAccessCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	hasAccessCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	hasAccessCmd.Flags().StringVar(&userId, "id", "", "user id")
	hasAccessCmd.Flags().BoolVar(&describeResource, "describe", false, "show extended info")

	hasAccessCmd.MarkFlagRequired("id")
}