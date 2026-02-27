/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/

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

var (
	email string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new shared mailbox",
	Long: `Use this command to create a new shared mailbox.
"ya360_admin:mail_write_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("mailbox add called")

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

		item := struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Email       string `json:"email"`
		}{
			Name:        name,
			Description: description,
			Email:       email,
		}

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/shared", helper.BaseUrl, orgId)
		payload, _ := json.Marshal(item)

		resp, err := helper.MakeRequest(url, "PUT", token, payload)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
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
	addCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	addCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	addCmd.Flags().StringVar(&name, "name", "", "shared mailbox name")
	addCmd.Flags().StringVar(&description, "description", "", "shared mailbox description")
	addCmd.Flags().StringVar(&email, "email", "", "shared mailbox email address")

	addCmd.MarkFlagRequired("name")
	addCmd.MarkFlagRequired("description")
	addCmd.MarkFlagRequired("email")
}
