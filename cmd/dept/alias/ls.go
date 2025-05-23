/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package alias

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

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Get a list of the department aliases",
	Long: `Use this command to get a list of the department aliases.
"directory:read_departments" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("dept alias ls called")

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
		t.AppendRow(table.Row{"Name", data.Name})
		t.AppendRow(table.Row{"Aliases", strings.Join(data.Aliases, ",")})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	lsCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	lsCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	lsCmd.Flags().IntVar(&deptId, "id", 0, "department id")

	lsCmd.MarkFlagRequired("id")
}
