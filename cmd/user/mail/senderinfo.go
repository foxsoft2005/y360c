/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package mail

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

// senderinfoCmd represents the senderinfo command
var senderinfoCmd = &cobra.Command{
	Use:   "sender-info",
	Short: "Show a sender info for the user",
	Long: `Use this command to show a sender info for the selected user.
"ya360_admin:mail_read_user_settings" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user mail sender-info called")

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

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
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
	senderinfoCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	senderinfoCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	senderinfoCmd.Flags().StringVar(&userId, "id", "", "user id")

	senderinfoCmd.MarkFlagRequired("id")

}
