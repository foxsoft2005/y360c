/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/
package member

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

// membersCmd represents the members command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a member to the group",
	Long: `Use this command to add a member to the selected group.
"directory:write_groups" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("group member add called")

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
			Id         string          `json:"id"`
			MemberType groupMemberType `json:"type"`
		}{
			Id:         memberId,
			MemberType: memberType,
		}

		var url = fmt.Sprintf("%s/directory/v1/org/%d/groups/%d/members", helper.BaseUrl, orgId, groupId)
		payload, _ := json.Marshal(item)

		resp, err := helper.MakeRequest(url, "POST", token, payload)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		var data model.GroupMemberResponse
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.Id})
		t.AppendRow(table.Row{"Type", data.Type})
		t.AppendRow(table.Row{"Added", data.Added})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	memberType = userMember

	addCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	addCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	addCmd.Flags().IntVar(&groupId, "id", 0, "group id")
	addCmd.Flags().Var(&memberType, "member-type", "member (user, group or department) to be added")
	addCmd.Flags().StringVar(&memberId, "member-id", "", "member id to be added")

	addCmd.MarkFlagRequired("id")
	addCmd.MarkFlagRequired("member-type")
	addCmd.MarkFlagRequired("member-id")
}
