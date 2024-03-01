/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package member

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// membersCmd represents the members command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "adds a member to the group",
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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/groups/%d/members", helper.BaseUrl, orgId, groupId)
		var payload = []byte(fmt.Sprintf(`{"id":"%s", "type":"%s"}`, memberId, memberType))

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

	addCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "Organization id")
	addCmd.Flags().StringVarP(&token, "token", "t", "", "Access token")
	addCmd.Flags().IntVar(&groupId, "id", 0, "Group id")
	addCmd.Flags().Var(&memberType, "memberType", "Member (user, group or department) to be added")
	addCmd.Flags().StringVar(&memberId, "memberId", "", "Member id to be added")

	addCmd.MarkFlagRequired("id")
	addCmd.MarkFlagRequired("memberType")
	addCmd.MarkFlagRequired("memberId")
}
