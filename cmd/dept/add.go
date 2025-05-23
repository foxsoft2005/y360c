/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package dept

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-json"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// membersCmd represents the members command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Create a new department",
	Long: `Use this command to create a new department.
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
			Name        *string `json:"name"`
			Label       *string `json:"label,omitempty"`
			Description *string `json:"description,omitempty"`
			ParentId    int     `json:"parentId"`
			ExternalId  *string `json:"externalId,omitempty"`
			HeadId      *string `json:"headId,omitempty"`
		}{
			Name:        helper.ToNullableString(name),
			Label:       helper.ToNullableString(label),
			Description: helper.ToNullableString(description),
			ParentId:    parentId,
			ExternalId:  helper.ToNullableString(externalId),
			HeadId:      helper.ToNullableString(headId),
		}

		var url = fmt.Sprintf("%s/directory/v1/org/%d/departments", helper.BaseUrl, orgId)
		payload, _ := json.Marshal(item)

		resp, err := helper.MakeRequest(url, "POST", token, payload)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
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
	addCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	addCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	addCmd.Flags().StringVar(&name, "name", "", "department name")
	addCmd.Flags().StringVar(&label, "label", "", "department label")
	addCmd.Flags().StringVar(&description, "description", "", "department description")
	addCmd.Flags().IntVar(&parentId, "parent-id", 1, "parent department id")
	addCmd.Flags().StringVar(&externalId, "external-id", "", "external department id")
	addCmd.Flags().StringVar(&headId, "head-id", "", "department head id")

	addCmd.MarkFlagRequired("name")
}
