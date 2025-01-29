/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package alias

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

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Get a list of the user mailbox aliases",
	Long: `Use this command to get a list of the user mailbox aliases.
"directory:read_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user alias ls called")

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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/users/%s", helper.BaseUrl, orgId, userId)

		resp, err := helper.MakeRequest(url, "GET", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		var data model.User
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.Id})
		t.AppendRow(table.Row{"Name", strings.TrimSpace(fmt.Sprintf("%s %s %s", data.Name.Last, data.Name.First, data.Name.Middle))})
		t.AppendRow(table.Row{"Aliases", strings.Join(data.Aliases, ",")})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	lsCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	lsCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	lsCmd.Flags().StringVar(&userId, "id", "", "user id")

	lsCmd.MarkFlagRequired("id")
}
