// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package mfa

import (
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-json"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get a 2FA status for domain users",
	Long: `Get a two-factor authentication (2FA) status for domain users.
"ya360_security:domain_2fa_write" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("mfa status called")

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

		resp, err := helper.MakeRequest(url, "GET", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
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
	statusCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	statusCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
}
