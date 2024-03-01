/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package disk

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var maxRec int

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "gets a list of events in the organization's Disk audit log",
	Long: `Use this command to retrieve a list of events in the organization's Disk audit log.
"ya360_security:audit_log_disk" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("disk log called")

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

		var url = fmt.Sprintf("%s/security/v1/org/%d/audit_log/disk?pageSize=%d", helper.BaseUrl, orgId, maxRec)

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

		var data model.DiskAuditLog
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Client IP", "Date", "Event Type", "Owner Login", "Path", "Resource File Id", "Size"})
		for _, e := range data.Events {
			created, _ := time.Parse(time.RFC3339, e.Date)
			t.AppendRow(table.Row{e.ClientIp, created, e.EventType, e.OwnerLogin, e.Path, e.ResourceFileId, e.Size})
		}

		t.AppendSeparator()
		t.Render()
	},
}

func init() {
	logCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "Organization id")
	logCmd.Flags().StringVarP(&token, "token", "t", "", "Access token")
	logCmd.Flags().IntVar(&maxRec, "max", 100, "Max records to retrieve")
}
