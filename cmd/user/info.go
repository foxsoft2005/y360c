/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package user

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Gets a user information",
	Long: `Use this command to get a user information by id.
"directory:read_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		if !asRaw {
			log.Println("user info called")
		}

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

		data, err := helper.GetUserById(orgId, token, userId)
		if err != nil {
			log.Fatalln("Unable to get user:", err)
		}

		if asRaw {
			buff, _ := json.MarshalIndent(data, "", "     ")
			fmt.Print(string(buff))
		} else {
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
		}
	},
}

func init() {
	infoCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	infoCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	infoCmd.Flags().StringVar(&userId, "id", "", "user id")
	infoCmd.Flags().BoolVar(&asRaw, "raw", false, "export as raw data")

	infoCmd.MarkFlagsOneRequired("id")
}
