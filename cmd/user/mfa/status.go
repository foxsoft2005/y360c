/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/
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
	Short: "Show a status of 2fa for the user",
	Long: `Use this command to show a status of two-factor auth (2fa) for the selected user.
"directory:read_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user mfa status called")

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

		if userEmail != "" {
			us, err := helper.GetUserByEmail(orgId, token, userEmail)
			if err != nil {
				log.Fatalln("Failed to get user by email", err)
			}

			if us == nil {
				log.Fatalf("User (mailbox) %s does not found", userEmail)
			}

			userId = us.Id
		}

		var url = fmt.Sprintf("%s/directory/v1/org/%d/users/%s/2fa", helper.BaseUrl, orgId, userId)

		resp, err := helper.MakeRequest(url, "GET", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		var data model.UserMfaSetup
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"User Id", data.UserId})
		t.AppendRow(table.Row{"Has 2FA", data.Has2fa})
		t.AppendRow(table.Row{"Has Phone", data.HasSecurityPhone})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	statusCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	statusCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	statusCmd.Flags().StringVar(&userId, "id", "", "user id")
	statusCmd.Flags().StringVar(&userEmail, "email", "", "user email address")

	statusCmd.MarkFlagsOneRequired("id", "email")
	statusCmd.MarkFlagsMutuallyExclusive("id", "email")
}
