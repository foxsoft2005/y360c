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
	toId         string
	notify       bool
	deleteEntry  bool
	sendAs       bool
	sendOnBehalf bool
)

// delegateCmd represents the delegate command
var delegateCmd = &cobra.Command{
	Use:   "delegate",
	Short: "delegates access to the mailbox to other user",
	Long: `Use this command to delegate access to the mailbox to other user.
"ya360_admin:mail_read_shared_mailbox_inventory" and "ya360_admin:mail_write_shared_mailbox_inventory"
permissions are required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user mail delegate called")

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

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mail/delegated?resourceId=%s&actorId=%s&sendMessages=%t", helper.BaseUrl, orgId, userId, toId, notify)

		var method = "POST"
		var entry model.MailAccessSettings

		if deleteEntry {
			method = "DELETE"
		} else {
			entry.Items = append(entry.Items, "imap_full_access")
			if sendAs {
				entry.Items = append(entry.Items, "send_as")
			}
			if sendOnBehalf {
				entry.Items = append(entry.Items, "send_on_behalf")
			}
		}

		payload, _ := json.Marshal(entry)
		resp, err := helper.MakeRequest(url, method, token, payload)
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

		var data model.MailAccessResponse
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Task Id", data.TaskId})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()

	},
}

func init() {
	delegateCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	delegateCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	delegateCmd.Flags().StringVar(&userId, "id", "", "user id")
	delegateCmd.Flags().StringVar(&toId, "toId", "", "user id to whom access is delegated")
	delegateCmd.Flags().BoolVar(&deleteEntry, "delete", false, "delete existing delegation")
	delegateCmd.Flags().BoolVar(&sendAs, "asUser", false, "make possible to send messages as user")
	delegateCmd.Flags().BoolVar(&sendOnBehalf, "onBehalf", false, "make possible to send messages on behalf of user")
	delegateCmd.Flags().BoolVar(&notify, "notify", false, "notify both users about delegation")

	delegateCmd.MarkFlagRequired("id")
	delegateCmd.MarkFlagRequired("toId")
	delegateCmd.MarkFlagsMutuallyExclusive("delete", "asUser")
	delegateCmd.MarkFlagsMutuallyExclusive("delete", "onBehalf")
}
