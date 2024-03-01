/*
Copyright © 2024 Kirill Chernetsky <kirill.chernetsky@linru.ru>
*/
package dept

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"linru.ru/y360c/helper"
	"linru.ru/y360c/model"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "gets department info",
	Long: `Use this command to get department information by id.
"directory:read_departments" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user dept info called")

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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/departments/%d", helper.BaseUrl, orgId, deptId)

		resp, err := helper.MakeRequest(url, "GET", token, nil)
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
		t.AppendRow(table.Row{"Label", data.Label})
		t.AppendRow(table.Row{"Email", data.Email})
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
	infoCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "Organization id")
	infoCmd.Flags().StringVarP(&token, "token", "t", "", "Access token")
	infoCmd.Flags().IntVar(&deptId, "id", 0, "Department id")

	infoCmd.MarkFlagRequired("id")
}
