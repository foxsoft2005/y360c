// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package mailbox

import (
	"log"
	"os"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	taskId string
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check status of the delegation task",
	Long: `Use this command to check status of the task which was created when delegation requested.
"ya360_admin:mail_write_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
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

		data, err := helper.CheckTaskById(orgId, token, taskId)
		if err != nil {
			log.Fatalln("Unable to get task:", err)
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
	statusCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	statusCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	statusCmd.Flags().StringVar(&taskId, "task-id", "", "task id")

	err := statusCmd.MarkFlagRequired("task-id")
	if err != nil {
		log.Fatalln("Error marking flag as required:", err)
	}
}
