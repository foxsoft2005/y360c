/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package mfa

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// disableCmd represents the disable command
var disableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disables mandatory 2FA for domain users",
	Long: `Use this command to disable madatory two-factor authentication (2FA) for the selected organization.
"ya360_security:domain_2fa_write" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mfa disable called")

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

		var url = fmt.Sprintf("%s/security/v1/org/%d/domain_2fa", helper.BaseUrl, orgId)

		resp, err := helper.MakeRequest(url, "DELETE", token, nil)
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

		var data model.MfaSetup
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Duration", data.Duration})
		t.AppendRow(table.Row{"Enabled", data.Enabled})
		t.AppendRow(table.Row{"Enabled At", data.EnabledAt})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	disableCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	disableCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
}
