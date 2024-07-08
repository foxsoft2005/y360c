/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package mail

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
	describeActor bool
)

// sharedWithCmd represents the sharedWith command
var sharedWithCmd = &cobra.Command{
	Use:   "sharedWith",
	Short: "gets all resources that have access to mailbox",
	Long: `Use this command to get all resources (users, groups) that have access to selected mailbox.
"ya360_admin:mail_read_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user mail sharedWith called")

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

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mail/delegated/%s/actors", helper.BaseUrl, orgId, userId)

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

		var data model.ActorList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		if describeActor {
			t.AppendHeader(table.Row{"Id", "Name", "Email", "Rights"})
			for _, e := range data.Items {
				user, _ := helper.GetUserById(orgId, token, e.ActorId)
				t.AppendRow(table.Row{e.ActorId, fmt.Sprintf("%s %s", user.Name.First, user.Name.Last), strings.Join(e.Items, ",")})
			}
		} else {
			t.AppendHeader(table.Row{"Id", "Rights"})
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
	sharedWithCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	sharedWithCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	sharedWithCmd.Flags().StringVar(&userId, "id", "", "user id")
	sharedWithCmd.Flags().BoolVar(&describeActor, "describe", false, "show extended info")

	sharedWithCmd.MarkFlagRequired("id")
}
