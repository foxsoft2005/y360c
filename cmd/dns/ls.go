/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package dns

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

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "gets a list of all DNS records for the domain",
	Long: `Use this command to retrieve a list of all DNS records which are set for the domain.
"directory:manage_dns" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("dns ls called")

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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/domains/%s/dns?perPage=%d", helper.BaseUrl, orgId, domain, maxRec)

		resp, err := helper.MakeRequest(url, "GET", token, nil)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if resp.HttpCode != 200 {
			var errorData model.ErrorResponse
			if err := json.Unmarshal(resp.Body, &errorData); err != nil {
				log.Fatalln("Unable to evaluate data:", err)
			}
			log.Fatalf("Response (HTTP %d): [%d] %s", resp.HttpCode, errorData.Code, errorData.Message)
		}

		var data model.DnsRecordList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"record id", "type", "name", "address", "ttl"})
		for _, e := range data.Records {
			t.AppendRow(table.Row{e.RecordId, e.Type, e.Name, e.Address, e.Ttl})
		}

		t.AppendSeparator()
		t.Render()

	},
}

func init() {
	lsCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "Organization id")
	lsCmd.Flags().StringVarP(&token, "token", "t", "", "Access token")
	lsCmd.Flags().IntVar(&maxRec, "max", 50, "Max records to retrieve")
	lsCmd.Flags().StringVar(&domain, "domain", "", "Domain name (example.com)")

	lsCmd.MarkFlagRequired("domain")
}
