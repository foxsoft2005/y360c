/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package domain

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get a Y360 domain information",
	Long: `Use this command to retrieve a Y360 domain by name.
"directory:read_domains" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("domain info called")

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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/domains?perPage=%d", helper.BaseUrl, orgId, maxRec)

		resp, err := helper.MakeRequest(url, "GET", token, nil)
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

		var data model.DomainList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		for _, e := range data.Domains {
			if strings.EqualFold(e.Name, domain) {
				t.AppendRow(table.Row{"Name", e.Name})
				t.AppendRow(table.Row{"Country", e.Country})
				t.AppendRow(table.Row{"Is master", e.Master})
				t.AppendRow(table.Row{"Is mx", e.MX})
				t.AppendRow(table.Row{"Delegated", e.Delegated})
				t.AppendRow(table.Row{"Verified", e.Verified})
				t.AppendRow(table.Row{"MX", e.Status.MX.Value})
				t.AppendRow(table.Row{"DKIM", e.Status.DKIM.Value})
				t.AppendRow(table.Row{"NS", e.Status.NS.Value})
				t.AppendRow(table.Row{"SPF", e.Status.SPF.Value})

				break
			}
		}
		t.AppendSeparator()
		t.SetAllowedRowLength(80)
		t.Style().Options.SeparateRows = true
		t.Render()

	},
}

func init() {
	infoCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	infoCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	infoCmd.Flags().IntVar(&maxRec, "max", 10, "max records to retrieve")
	infoCmd.Flags().StringVar(&domain, "domain", "", "domain name (example.com)")

	infoCmd.MarkFlagRequired("domain")
}
