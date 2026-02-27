/*
Copyright © 2024 Kirill Chernetsky aka foxsoft2005
*/
package user

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

var force bool

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove the user",
	Long: `Use this command to completely (all incl. emails, files, messages, etc.) remove the selected user.
"directory:write_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("user rm is called")

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

		if !force {
			if !helper.Confirm("Do you REALLY want to DELETE the selected entity (y[es]|no)?") {
				log.Fatal("Aborted by the user")
			}
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

		user, err := helper.GetUserById(orgId, token, userId)
		if err != nil {
			log.Fatalln("Unable to get user:", err)
		}

		if user.IsEnabled {
			log.Fatalf("User %s (%s) is enabled and cannot be deleted", user.Id, user.Email)
		}

		var url = fmt.Sprintf("%s/directory/v1/org/%d/users/%s", helper.BaseUrl, orgId, userId)

		resp, err := helper.MakeRequest(url, "DELETE", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		var data model.UserDeletionResponse
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"User Id", data.UserId})
		t.AppendRow(table.Row{"Deleted", data.Deleted})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	rmCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	rmCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	rmCmd.Flags().StringVar(&userId, "id", "", "user id to remove")
	rmCmd.Flags().StringVar(&userEmail, "email", "", "user email address")

	rmCmd.Flags().BoolVar(&force, "force", false, "force deletion")

	rmCmd.MarkFlagsOneRequired("id", "email")
	rmCmd.MarkFlagsMutuallyExclusive("id", "email")
}
