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
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Removes a member from the group",
	Long: `Use this command to remove a member from the selected group.
"directory:write_groups" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("group member rm called")

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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/groups/%d/members/%s/%s", helper.BaseUrl, orgId, groupId, memberType, memberId)

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

		var data model.GroupMemberResponse
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.Id})
		t.AppendRow(table.Row{"Type", data.Type})
		t.AppendRow(table.Row{"Deleted", data.Deleted})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	memberType = userMember

	rmCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	rmCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	rmCmd.Flags().IntVar(&groupId, "id", 0, "group id")
	rmCmd.Flags().Var(&memberType, "memberType", "member type to be deleted")
	rmCmd.Flags().StringVar(&memberId, "memberId", "", "member id to be deleted")

	rmCmd.MarkFlagRequired("id")
	rmCmd.MarkFlagRequired("memberType")
	rmCmd.MarkFlagRequired("memberId")
}
