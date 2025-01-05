/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package auth

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

var authTtl int

// setCookieCmd represents the setCookie command
var setCookieTtlCmd = &cobra.Command{
	Use:   "setCookieTtl",
	Short: "Set a TTL of the user session cookies",
	Long: `Use this command to set a TTL (seconds) of the user session cookies.
"ya360_security:domain_sessions_write" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("auth setCookie called")

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

		var url = fmt.Sprintf("%s/security/v1/org/%d/domain_sessions", helper.BaseUrl, orgId)

		payload, _ := json.Marshal(model.CookieTTL{AuthTTL: authTtl})

		resp, err := helper.MakeRequest(url, "POST", token, payload)
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

		var data model.CookieTTL
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"New Auth TTL (sec.)", data.AuthTTL})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	setCookieTtlCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	setCookieTtlCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	setCookieTtlCmd.Flags().IntVar(&authTtl, "ttl", 0, "auth cookies termination timeout (sec.)")

	setCookieTtlCmd.MarkFlagRequired("ttl")
}
