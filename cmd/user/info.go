/*
Copyright © 2024 Kirill Chernetsky <kirill.chernetsky@linru.ru>
*/
package user

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
	Short: "gets a user information",
	Long: `Use this command to get a user information by id.
"directory:read_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user info called")

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

		if resp.HttpCode != 200 {
			var errorData model.ErrorResponse
			if err := json.Unmarshal(resp.Body, &errorData); err != nil {
				log.Fatalln("Unable to evaluate data:", err)
			}
			log.Fatalf("Response (HTTP %d): [%d] %s", resp.HttpCode, errorData.Code, errorData.Message)
		}

		var data model.User
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.Id})
		t.AppendRow(table.Row{"External Id", data.ExternalId})
		t.AppendRow(table.Row{"Name", strings.TrimSpace(fmt.Sprintf("%s %s %s", data.Name.Last, data.Name.First, data.Name.Middle))})
		t.AppendRow(table.Row{"Display Name", data.DisplayName})
		t.AppendRow(table.Row{"Nickname", data.Nickname})
		t.AppendRow(table.Row{"Email", data.Email})
		t.AppendRow(table.Row{"About", data.About})
		t.AppendRow(table.Row{"Gender", data.Gender})
		t.AppendRow(table.Row{"Birthday", data.Birthday})
		t.AppendRow(table.Row{"Department Id", data.DepartmentId})
		t.AppendRow(table.Row{"Position", data.Position})
		t.AppendRow(table.Row{"Language", data.Language})
		t.AppendRow(table.Row{"Timezone", data.Timezone})
		t.AppendRow(table.Row{"Aliases", strings.Join(data.Aliases, ",")})
		t.AppendRow(table.Row{"Groups", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(data.Groups)), ","), "[]")})
		t.AppendRow(table.Row{"Enabled", data.IsEnabled})
		t.AppendRow(table.Row{"Admin", data.IsAdmin})
		t.AppendRow(table.Row{"Dismissed", data.IsDismissed})
		t.AppendRow(table.Row{"Robot", data.IsRobot})
		t.AppendRow(table.Row{"Created At", data.CreatedAt})
		t.AppendRow(table.Row{"Updated At", data.UpdatedAt})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	infoCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "Organization id")
	infoCmd.Flags().StringVarP(&token, "token", "t", "", "Access token")
	infoCmd.Flags().StringVar(&userId, "id", "", "User id")

	infoCmd.MarkFlagsOneRequired("id")
}
