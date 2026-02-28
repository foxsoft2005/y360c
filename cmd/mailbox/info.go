// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package mailbox

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-json"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get info about the shared mailbox",
	Long: `Use this command to retrieve information about the shared mailbox.
"ya360_admin:mail_read_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("mailbox info called")

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

			if us == nil {
				log.Fatalf("User (mailbox) %s does not found", mailboxName)
			}

			mailboxId = us.Id
		}

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/shared/%s", helper.BaseUrl, orgId, mailboxId)

		resp, err := helper.MakeRequest(url, "GET", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		var data model.MailboxInfo
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.Id})
		t.AppendRow(table.Row{"Email", data.Email})
		t.AppendRow(table.Row{"Name", data.Name})
		t.AppendRow(table.Row{"Description", data.Description})
		t.AppendRow(table.Row{"Created At", data.CreatedAt})
		t.AppendRow(table.Row{"Updated At", data.UpdatedAt})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()

	},
}

func init() {
	infoCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	infoCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	infoCmd.Flags().StringVar(&mailboxId, "id", "", "shared mailbox id")
	infoCmd.Flags().StringVar(&mailboxName, "email", "", "mailbox (or user) email address")

	infoCmd.MarkFlagsOneRequired("id", "email")
	infoCmd.MarkFlagsMutuallyExclusive("id", "email")
}
