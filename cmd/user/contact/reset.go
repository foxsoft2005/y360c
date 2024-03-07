/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package contact

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

// rmCmd represents the rm command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "resets a contact information",
	Long: `Use this command to remove ALL manually entered contact information for the selected user.
"directory:write_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user contact rm called")

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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/users/%s/contacts", helper.BaseUrl, orgId, userId)

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

		var data model.User
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.Id, ""})
		t.AppendRow(table.Row{"Name", strings.TrimSpace(fmt.Sprintf("%s %s %s", data.Name.Last, data.Name.First, data.Name.Middle)), ""})
		for index, item := range data.Contacts {
			if index == 0 {
				t.AppendRow(table.Row{"Contact info", "", ""})
			}

			var s = item.Type
			if item.Main {
				s = fmt.Sprintf("%s (main)", item.Type)
			}

			var s1 = item.Value
			if item.Synthetic {
				s1 = fmt.Sprintf("%s (readonly)", item.Value)
			}

			t.AppendRow(table.Row{"", s, s1})
		}
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	resetCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	resetCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	resetCmd.Flags().StringVar(&userId, "id", "", "user id")

	resetCmd.MarkFlagRequired("id")
}
