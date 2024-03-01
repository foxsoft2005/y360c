/*
Copyright © 2024 Kirill Chernetsky <kirill.chernetsky@linru.ru>
*/
package alias

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

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "removes an alias of the user mailbox",
	Long: `Use this command to remove an alias for the user mailbox.
"directory:write_departments" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("dept alias rm called")

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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/departments/%d/aliases/%s", helper.BaseUrl, orgId, deptId, alias)

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

		var data model.RmAliasResponse
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Alias", "Removed"})
		t.AppendRow(table.Row{data.Alias, data.Removed})
		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	rmCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "Organization id")
	rmCmd.Flags().StringVarP(&token, "token", "t", "", "Access token")
	rmCmd.Flags().IntVar(&deptId, "id", 0, "Department id")
	rmCmd.Flags().StringVar(&alias, "alias", "", "Alias to be deleted")

	rmCmd.MarkFlagRequired("id")
	rmCmd.MarkFlagRequired("alias")
}