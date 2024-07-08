/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package mail

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

var (
	taskId string
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "checks status of the delegation task",
	Long: `Use this command to check status of the task which was created when delegation requested.
"ya360_admin:mail_read_shared_mailbox_inventory" and "ya360_admin:mail_write_shared_mailbox_inventory"
permissions are required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user mail status called")

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

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mail/delegated/tasks/%s", helper.BaseUrl, orgId, taskId)

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

		var data model.TaskStatusResponse
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Status", data.Status})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()

	},
}

func init() {
	statusCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	statusCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	statusCmd.Flags().StringVar(&taskId, "taskId", "", "task id")

	statusCmd.MarkFlagRequired("taskId")
}
