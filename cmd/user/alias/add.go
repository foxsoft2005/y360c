/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "adds an alias for the user mailbox",
	Long: `Use this command to add an alias for the user mailbox.
"directory:write_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user alias add called")

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

		var payload = []byte(fmt.Sprintf(`{"alias":"%s"}`, alias))
		var url = fmt.Sprintf("%s/directory/v1/org/%d/users/%s/aliases", helper.BaseUrl, orgId, userId)

		resp, err := helper.MakeRequest(url, "POST", token, payload)
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
	addCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	addCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	addCmd.Flags().StringVar(&userId, "id", "", "user id")
	addCmd.Flags().StringVar(&alias, "alias", "", "mailbox alias")

	addCmd.MarkFlagRequired("id")
	addCmd.MarkFlagRequired("alias")
}
