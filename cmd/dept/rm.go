/*
Copyright © 2024 Kirill Chernetsky <kirill.chernetsky@linru.ru>
*/
package dept

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"linru.ru/y360c/helper"
	"linru.ru/y360c/model"
)

// membersCmd represents the members command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "removes a department",
	Long: `Use this command to remove a department.
"directory:write_departments" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("dept rm called")

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

		resp, err := helper.MakeRequest(url, "DELETE", token, nil)
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
		t.AppendRow(table.Row{"Removed", data.Removed})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	rmCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "Organization id")
	rmCmd.Flags().StringVarP(&token, "token", "t", "", "Access token")
	rmCmd.Flags().IntVar(&deptId, "id", 0, "Department id")

	rmCmd.MarkFlagRequired("id")
}
