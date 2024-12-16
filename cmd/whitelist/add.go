/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package whitelist

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

// lsCmd represents the ls command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds the list of the allowed IPs and/or CIDRs",
	Long: `Use this command to add the list of the allowed IP addresses and/or CIDRs.
"ya360_admin:mail_write_antispam_settings" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("whitelist add called")

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

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mail/antispam/allowlist/ips", helper.BaseUrl, orgId)
		payload, _ := json.Marshal(model.WhiteList{AllowList: allowed})

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

		var data model.WhiteList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Allow list"})
		for _, e := range data.AllowList {
			t.AppendRow(
				table.Row{e},
			)

		}
		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	addCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	addCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	addCmd.Flags().StringArrayVar(&allowed, "allowed", nil, `allowed IP-address ("77.88.21.249") or CIDR ("77.88.54.0/23")`)

	addCmd.MarkFlagRequired("allowed")
}
