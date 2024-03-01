/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package group

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

// membersCmd represents the members command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "creates a new group",
	Long: `Use this command to create a new group.
"directory:write_groups" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("group add called")

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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/groups", helper.BaseUrl, orgId)
		payload, _ := json.Marshal(model.Group{Name: name, Description: description, Label: label, ExternalId: externalId, AdminIds: admins})

		resp, err := helper.MakeRequest(url, "POST", token, payload)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if resp.HttpCode != 200 {
			var errorData model.ErrorResponse
			if err := json.Unmarshal(resp.Body, &errorData); err != nil {
				log.Fatalln("Unable to evaluate data:", err)
			}
			log.Fatalf("Response (HTTP %d): [%d] %s", resp.HttpCode, errorData.Code, errorData.Message)
		}

		var data model.Group
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.Id})
		t.AppendRow(table.Row{"Type", data.Type})
		t.AppendRow(table.Row{"External Id", data.ExternalId})
		t.AppendRow(table.Row{"Name", data.Name})
		t.AppendRow(table.Row{"Description", data.Description})
		t.AppendRow(table.Row{"Email", data.Email})
		t.AppendRow(table.Row{"Label", data.Label})
		t.AppendRow(table.Row{"Author Id", data.AuthorId})
		t.AppendRow(table.Row{"Admin Ids", strings.Join(data.AdminIds, ",")})
		t.AppendRow(table.Row{"Aliases", strings.Join(data.Aliases, ",")})
		t.AppendRow(table.Row{"Member Of", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(data.MemberOf)), ","), "[]")})
		t.AppendRow(table.Row{"Members Count", data.MembersCount})
		t.AppendRow(table.Row{"Created At", data.CreatedAt})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	addCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "Organization id")
	addCmd.Flags().StringVarP(&token, "token", "t", "", "Access token")
	addCmd.Flags().StringVar(&name, "name", "", "Group name")
	addCmd.Flags().StringVar(&label, "label", "", "Group label")
	addCmd.Flags().StringVar(&description, "description", "", "Group description")
	addCmd.Flags().StringVar(&externalId, "externalId", "", "External id")
	addCmd.Flags().StringArrayVar(&admins, "admin", nil, "Group administrators ids")

	addCmd.MarkFlagRequired("name")
}
