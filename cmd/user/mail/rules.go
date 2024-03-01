/*
Copyright © 2024 Kirill Chernetsky <kirill.chernetsky@linru.ru>
*/
package mail

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

// rulesCmd represents the rules command
var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "shows mailbox rules for the user",
	Long: `Use this command to show mailbox rules (autoreplies, forwards) for the selected user.
"ya360_admin:mail_read_user_settings" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user mail rules called")

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

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mail/users/%s/settings/user_rules", helper.BaseUrl, orgId, userId)

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

		var data model.UserMailRules
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Autoreplies", "", ""})
		for _, el := range data.Autoreplies {
			t.AppendRow(table.Row{"", "Rule Id", el.RuleId})
			t.AppendRow(table.Row{"", "Rule Name", el.RuleName})
			t.AppendRow(table.Row{"", "Text", el.Text})
		}
		t.AppendRow(table.Row{"Forwards", "", ""})
		for _, el := range data.Forwards {
			t.AppendRow(table.Row{"", "Rule Id", el.RuleId})
			t.AppendRow(table.Row{"", "Rule Name", el.RuleName})
			t.AppendRow(table.Row{"", "Address", el.Address})
			t.AppendRow(table.Row{"", "With Store", el.WithStore})
		}
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()

	},
}

func init() {
	rulesCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "Organization id")
	rulesCmd.Flags().StringVarP(&token, "token", "t", "", "Access token")
	rulesCmd.Flags().StringVar(&userId, "id", "", "User id")

	rulesCmd.MarkFlagRequired("id")

}
