/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package auth

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/spf13/cobra"
)

var userId string

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logs out a user on all devices",
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

		var url = fmt.Sprintf("%s/security/v1/org/%d/domain_sessions/users/%s/logout", helper.BaseUrl, orgId, userId)

		resp, err := helper.MakeRequest(url, "PUT", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if resp.HttpCode != 200 {
			var errData model.ErrorResponse
			if err := json.Unmarshal(resp.Body, &errData); err != nil {
				log.Fatalln("Unable to evaluate error:", err)
			}
			log.Fatalf("http %d: [%d] %s", resp.HttpCode, errData.Code, errData.Message)
		}

		log.Printf("User %s was logged out successfully", userId)
	},
}

func init() {
	logoutCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	logoutCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	logoutCmd.Flags().StringVar(&userId, "userId", "", "user id to log out")

	logoutCmd.MarkFlagRequired("userId")
}
