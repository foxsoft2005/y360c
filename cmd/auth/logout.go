// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package auth

import (
	"fmt"
	"log"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/spf13/cobra"
)

var (
	userId    string
	userEmail string
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out a user on all devices",
	Long: `Use this command to log out the selected user on all devices.
"ya360_security:domain_sessions_write" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("auth logout called")

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

		if !helper.Confirm("Do you REALLY want to log the user out on ALL devices (y[es]|no)?") {
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

		var url = fmt.Sprintf("%s/security/v1/org/%d/domain_sessions/users/%s/logout", helper.BaseUrl, orgId, userId)

		resp, err := helper.MakeRequest(url, "PUT", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		log.Printf("User %s was logged out successfully", userId)
	},
}

func init() {
	logoutCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	logoutCmd.Flags().StringVarP(&token, "token", "t", "", "access token")

	logoutCmd.Flags().StringVar(&userId, "id", "", "id of the user to log out")
	logoutCmd.Flags().StringVar(&userEmail, "email", "", "email of the user to log out ")

	logoutCmd.MarkFlagsOneRequired("id", "email")
	logoutCmd.MarkFlagsMutuallyExclusive("id", "email")
}
