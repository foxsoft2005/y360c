/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package mfa

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets a phone number used for 2fa",
	Long: `Use this command to reset a phone number used for two-factor auth (2fa) by the user.
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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/users/%s/2fa", helper.BaseUrl, orgId, userId)

		resp, err := helper.MakeRequest(url, "DELETE", token, nil)
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

		log.Printf("2fa settings for the user %s were successfully cleared", userId)
	},
}

func init() {
	resetCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	resetCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	resetCmd.Flags().StringVar(&userId, "id", "", "user id")

	resetCmd.MarkFlagRequired("id")
}
