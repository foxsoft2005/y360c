/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package admin

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
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove all admins from the group",
	Long: `Use this command to remove all admins from the selected group.
"directory:write_groups" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("group admin rm called")

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

		if !helper.Confirm("Do you REALLY want to DELETE the selected entity (y[es]|no)?") {
			log.Fatal("Aborted, exiting")
		}

		var url = fmt.Sprintf("%s/directory/v1/org/%d/groups/%d/admins", helper.BaseUrl, orgId, groupId)

		resp, err := helper.MakeRequest(url, "DELETE", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
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
	rmCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	rmCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	rmCmd.Flags().IntVar(&groupId, "id", 0, "group id")

	rmCmd.MarkFlagRequired("id")
}
