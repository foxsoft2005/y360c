/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package org

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

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "gets a list of the organizations",
	Long: `Use this command to retrieve a list of Y360 organizations.
"directory:read_organization" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("org ls called")

		if token == "" {
			t, err := helper.GetToken()
			if err != nil {
				log.Fatalln("Incorrect settings:", err)
			}
			token = t
		}

		var url = fmt.Sprintf("%s/directory/v1/org", helper.BaseUrl)

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

		var data model.OrganizationList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"id", "name", "phone", "fax", "email", "subscription"})
		for _, e := range data.Organizations {
			t.AppendRow(table.Row{e.Id, e.Name, e.Phone, e.Fax, e.Email, e.SubsciptionPlan})
		}
		t.AppendSeparator()
		t.Render()

	},
}

func init() {
	listCmd.Flags().StringVarP(&token, "token", "t", "", "Access token")
	listCmd.Flags().IntVarP(&maxRec, "max", "m", 100, "Max records to retrieve")

}
