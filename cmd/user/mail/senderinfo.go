/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package mail

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

// senderinfoCmd represents the senderinfo command
var senderinfoCmd = &cobra.Command{
	Use:   "senderinfo",
	Short: "shows a sender info for the user",
	Long: `Use this command to show a sender info for the selected user.
"ya360_admin:mail_read_user_settings" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user mail senderinfo called")

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

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mail/users/%s/settings/sender_info", helper.BaseUrl, orgId, userId)

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

		var data model.UserSenderInfo
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Default From", "", data.DefaultFrom})
		t.AppendRow(table.Row{"From Name", "", data.FromName})
		t.AppendRow(table.Row{"Sign Position", "", data.SignPosition})
		for idx, el := range data.Signs {
			t.AppendRow(table.Row{fmt.Sprintf("Sign %d", idx+1), "", ""})
			t.AppendRow(table.Row{"", "Emails", strings.Join(el.Emails, ",")})
			t.AppendRow(table.Row{"", "Is Default", el.IsDefault})
			t.AppendRow(table.Row{"", "Lang", el.Lang})
			t.AppendRow(table.Row{"", "Text", el.Text})
		}
		t.AppendSeparator()
		t.SetAllowedRowLength(100)
		t.Style().Options.SeparateRows = true
		t.Render()

	},
}

func init() {
	senderinfoCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "Organization id")
	senderinfoCmd.Flags().StringVarP(&token, "token", "t", "", "Access token")
	senderinfoCmd.Flags().StringVar(&userId, "id", "", "User id")

	senderinfoCmd.MarkFlagRequired("id")

}
