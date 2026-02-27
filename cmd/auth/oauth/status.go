/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/
package oauth

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

var (
	activate   bool
	deactivate bool
)

// command definition
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Manage OAuth status",
	Long: `Use this command to enable, disable or show status of OAuth for external services.
"ya360_security:domain_settings_write" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("auth oauth status called")

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

		var url = fmt.Sprintf("%s/security/v1/org/%d/oauth_access_restriction", helper.BaseUrl, orgId)

		var method = "GET"

		if activate {
			method = "POST"
		}

		if deactivate {
			method = "DELETE"
		}

		resp, err := helper.MakeRequest(url, method, token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		if !activate && !deactivate {
			var data model.OAuthStatusResponse

			if err := json.Unmarshal(resp.Body, &data); err != nil {
				log.Fatalln("Unable to evaluate data:", err)
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendRow(table.Row{"Restricted", data.Restricted})
			t.AppendSeparator()
			t.Style().Options.SeparateRows = true
			t.Render()
		} else {
			log.Println("Successfully done")
		}
	},
}

func init() {
	statusCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	statusCmd.Flags().StringVarP(&token, "token", "t", "", "access token")

	statusCmd.Flags().BoolVar(&activate, "enable", false, "enable OAuth")
	statusCmd.Flags().BoolVar(&deactivate, "disable", false, "disable OAuth")

	statusCmd.MarkFlagsMutuallyExclusive("enable", "disable")
}
