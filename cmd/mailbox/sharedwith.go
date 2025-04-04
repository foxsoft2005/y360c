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
	describeActor bool
)

// sharedWithCmd represents the sharedWith command
var sharedWithCmd = &cobra.Command{
	Use:   "shared-with",
	Short: "Get all resources that have access to mailbox",
	Long: `Use this command to get all resources (users, groups) that have access to selected mailbox.
"ya360_admin:mail_read_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user mail shared-with called")

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

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/actors/%s", helper.BaseUrl, orgId, mailboxId)

		resp, err := helper.MakeRequest(url, "GET", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		var data model.ActorList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		if describeActor {
			t.AppendHeader(table.Row{"Id", "Name", "Email", "Roles"})
			for _, e := range data.Items {
				user, _ := helper.GetUserById(orgId, token, e.ActorId)
				t.AppendRow(table.Row{e.ActorId, fmt.Sprintf("%s %s", user.Name.First, user.Name.Last), user.Email, strings.Join(e.Items, ",")})
			}
		} else {
			t.AppendHeader(table.Row{"Id", "Roles"})
			for _, e := range data.Items {
				t.AppendRow(table.Row{e.ActorId, strings.Join(e.Items, ",")})
			}
		}
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	sharedWithCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	sharedWithCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	sharedWithCmd.Flags().StringVar(&mailboxId, "id", "", "mailbox id")
	sharedWithCmd.Flags().StringVar(&mailboxName, "email", "", "mailbox (or user) email address")

	sharedWithCmd.Flags().BoolVar(&describeActor, "describe", false, "show extended info")

	sharedWithCmd.MarkFlagsOneRequired("id", "email")
	sharedWithCmd.MarkFlagsMutuallyExclusive("id", "email")
}
