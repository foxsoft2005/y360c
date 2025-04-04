/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package mailbox

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-json"

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
	Use:   "has-access",
	Short: "Get mailboxes that user has access to",
	Long: `Use this command to get mailboxes that selected user has access to.
"ya360_admin:mail_read_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("mailbox has-access called")

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

		if mailboxName != "" {
			us, err := helper.GetUserByEmail(orgId, token, mailboxName)
			if err != nil {
				log.Fatalln("Failed to get user by email", err)
			}

			mailboxId = us.Id
		}

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/resources/%s", helper.BaseUrl, orgId, mailboxId)

		resp, err := helper.MakeRequest(url, "GET", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
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
	hasAccessCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	hasAccessCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	hasAccessCmd.Flags().StringVar(&mailboxId, "id", "", "user id")
	hasAccessCmd.Flags().StringVar(&mailboxName, "email", "", "mailbox (or user) email address")

	hasAccessCmd.Flags().BoolVar(&describeResource, "describe", false, "show extended info")

	hasAccessCmd.MarkFlagsOneRequired("id", "email")
	hasAccessCmd.MarkFlagsMutuallyExclusive("id", "email")
}
