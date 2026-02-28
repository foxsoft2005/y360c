// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package mfa

import (
	"fmt"
	"log"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset 2fa settings",
	Long: `Use this command to reset two-factor auth (2fa) settings for the user.
"directory:write_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user mfa reset called")

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

		if !helper.Confirm("Do you REALLY want to RESET 2fa settings for the user (y[es]|no)?") {
			log.Fatalln("Aborted by the user")
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

		resp, err := helper.MakeRequest(url, "DELETE", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		log.Printf("2fa settings for the user %s were successfully cleared", userId)
	},
}

func init() {
	resetCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	resetCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	resetCmd.Flags().StringVar(&userId, "id", "", "user id")
	resetCmd.Flags().StringVar(&userEmail, "email", "", "user email address")

	resetCmd.MarkFlagsOneRequired("id", "email")
	resetCmd.MarkFlagsMutuallyExclusive("id", "email")
}
