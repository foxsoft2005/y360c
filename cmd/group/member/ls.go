/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package member

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
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Get a members of the group",
	Long: `Use this command to retrieve a members of the selected group.
"directory:read_groups" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("group member ls called")

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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/groups/%d/members", helper.BaseUrl, orgId, groupId)

		resp, err := helper.MakeRequest(url, "GET", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		var data model.GroupMemberList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		for i, e := range data.Departments {
			if i == 0 {
				t.AppendRow(table.Row{"Departments:", "", ""})
			}

			t.AppendRow(table.Row{"", "Id, Name", fmt.Sprintf("%d, %s", e.Id, e.Name)})
		}
		for i, e := range data.Groups {
			if i == 0 {
				t.AppendRow(table.Row{"Groups:", "", ""})
			}

			t.AppendRow(table.Row{"", "Id, Name", fmt.Sprintf("%d, %s", e.Id, e.Name)})
		}
		for i, e := range data.Users {
			if i == 0 {
				t.AppendRow(table.Row{"Users:", "", ""})
			}

			t.AppendRow(table.Row{"", "Id, Name", fmt.Sprintf("%s, %s", e.Id, strings.TrimSpace(fmt.Sprintf("%s %s %s", e.Name.Last, e.Name.First, e.Name.Middle)))})
		}
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	lsCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	lsCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	lsCmd.Flags().IntVar(&groupId, "id", 0, "group id")

	lsCmd.MarkFlagRequired("id")
}
