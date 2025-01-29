/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package user

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	about       string
	birthday    string
	displayName string
	externalId  string
	position    string
	isAdmin     helper.EnumYesNo
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set an attributes of the user",
	Long: `Use this command to set an attributes (name, birthday, department, etc.) of the selected user.
"directory:write_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("user set is called")

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

		if birthday != "" {
			_, err := time.Parse(time.DateOnly, birthday)
			if err != nil {
				log.Fatalln("Birthdate cannot be parsed:", err)
			}
		}

		item := struct {
			About       *string `json:"about,omitempty"`
			Birthday    *string `json:"birthday,omitempty"`
			DisplayName *string `json:"displayName,omitempty"`
			ExternalId  *string `json:"externalId,omitempty"`
			Position    *string `json:"position,omitempty"`
			IsAdmin     *bool   `json:"isAdmin,omitempty"`
		}{
			About:       helper.ToNullableString(about),
			Birthday:    helper.ToNullableString(birthday),
			DisplayName: helper.ToNullableString(displayName),
			ExternalId:  helper.ToNullableString(externalId),
			Position:    helper.ToNullableString(position),
			IsAdmin:     helper.EnumYesNoToBool(isAdmin),
		}

		var url = fmt.Sprintf("%s/directory/v1/org/%d/users/%s", helper.BaseUrl, orgId, userId)
		payload, _ := json.Marshal(item)

		resp, err := helper.MakeRequest(url, "PATCH", token, payload)
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
		t.AppendRow(table.Row{"About", data.About})
		t.AppendRow(table.Row{"Birthday", data.Birthday})
		t.AppendRow(table.Row{"Display Name", data.DisplayName})
		t.AppendRow(table.Row{"External Id", data.ExternalId})
		t.AppendRow(table.Row{"Position", data.Position})
		t.AppendRow(table.Row{"Is Admin", data.IsAdmin})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	setCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	setCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	setCmd.Flags().StringVar(&userId, "id", "", "user id")

	setCmd.Flags().StringVar(&about, "about", "", "about an user")
	setCmd.Flags().StringVar(&birthday, "birthday", "", "user bithday (YYYY-MM-DD)")
	setCmd.Flags().StringVar(&displayName, "display-name", "", "user display name")
	setCmd.Flags().StringVar(&externalId, "external-id", "", "user external id")
	setCmd.Flags().StringVar(&position, "position", "", "user position")
	setCmd.Flags().Var(&isAdmin, "is-admin", "user has admin permissions (yes or no)")

	setCmd.MarkFlagRequired("id")
	setCmd.MarkFlagsOneRequired("about", "birthday", "display-name", "external-id", "position", "is-admin")
}
