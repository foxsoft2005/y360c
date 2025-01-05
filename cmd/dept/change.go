/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package dept

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
var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "Changes the existing department",
	Long: `Use this command to change the existing department.
"directory:write_departments" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("dept add called")

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
			Name        string `json:"name,omitempty"`
			Label       string `json:"label,omitempty"`
			Description string `json:"description,omitempty"`
			ParentId    int    `json:"parentId,omitempty"`
			ExternalId  string `json:"externalId,omitempty"`
			HeadId      string `json:"headId,omitempty"`
		}{
			Name:        name,
			Label:       label,
			Description: description,
			ParentId:    parentId,
			ExternalId:  externalId,
			HeadId:      headId,
		}

		var url = fmt.Sprintf("%s/directory/v1/org/%d/departments/%d", helper.BaseUrl, orgId, deptId)
		payload, _ := json.Marshal(item)

		resp, err := helper.MakeRequest(url, "PATCH", token, payload)
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

		var data model.Department
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.Id})
		t.AppendRow(table.Row{"Parent Id", data.ParentId})
		t.AppendRow(table.Row{"External Id", data.ExternalId})
		t.AppendRow(table.Row{"Name", data.Name})
		t.AppendRow(table.Row{"Description", data.Description})
		t.AppendRow(table.Row{"Email", data.Email})
		t.AppendRow(table.Row{"Label", data.Label})
		t.AppendRow(table.Row{"Head Id", data.HeadId})
		t.AppendRow(table.Row{"Aliases", strings.Join(data.Aliases, ",")})
		t.AppendRow(table.Row{"Members Count", data.MembersCount})
		t.AppendRow(table.Row{"Created At", data.CreatedAt})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	changeCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	changeCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	changeCmd.Flags().IntVar(&deptId, "id", 0, "department id")
	changeCmd.Flags().StringVar(&name, "name", "", "department name")
	changeCmd.Flags().StringVar(&label, "label", "", "department label")
	changeCmd.Flags().StringVar(&description, "description", "", "department description")
	changeCmd.Flags().IntVar(&parentId, "parentId", 0, "parent department id")
	changeCmd.Flags().StringVar(&externalId, "externalId", "", "external department id")
	changeCmd.Flags().StringVar(&headId, "headId", "", "department head id")

	addCmd.MarkFlagRequired("id")
	addCmd.MarkFlagsOneRequired("name", "label", "description", "parentId", "externalId", "headId")
}
