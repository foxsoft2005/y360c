/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/
package mailbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type notificationType string

const notifyAll notificationType = "all"

func (e *notificationType) String() string {
	return string(*e)
}

func (e *notificationType) Set(v string) error {
	switch v {
	case "all", "delegates", "none":
		*e = notificationType(v)
		return nil
	default:
		return errors.New(`must be one of "all", or "delegates", or "none"`)
	}
}

func (e *notificationType) Type() string {
	return "notificationType"
}

var (
	toId                 string
	notify               notificationType
	mailboxOwner         bool
	mailboxImap          bool
	mailboxSender        bool
	mailboxLimitedSender bool
)

// setaccessCmd represents the delegate command
var setaccessCmd = &cobra.Command{
	Use:   "setAccess",
	Short: "Set access to the mailbox to other user",
	Long: `Use this command to set access to the mailbox to other user.
"ya360_admin:mail_write_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("mailbox setaccess called")

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

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/set/%s?actorId=%s&notify=%s", helper.BaseUrl, orgId, userId, toId, notify)

		var method = "POST"
		var entry model.MailAccessSettings

		if mailboxOwner {
			entry.Items = append(entry.Items, "shared_mailbox_owner")
		} else {
			if mailboxImap {
				entry.Items = append(entry.Items, "shared_mailbox_imap_admin")
			}

			if mailboxSender {
				entry.Items = append(entry.Items, "shared_mailbox_sender")
			}

			if mailboxLimitedSender {
				entry.Items = append(entry.Items, "shared_mailbox_half_sender")
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
	notify = notifyAll

	setaccessCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	setaccessCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	setaccessCmd.Flags().StringVar(&mailboxId, "id", "", "mailbox (or user) id")
	setaccessCmd.Flags().StringVar(&toId, "toId", "", "user id to whom access is delegated")
	setaccessCmd.Flags().BoolVar(&mailboxOwner, "owner", false, "full access to the mailbox (includes all other access types)")
	setaccessCmd.Flags().BoolVar(&mailboxImap, "readonly", false, "read-only access to the mailbox via IMAP")
	setaccessCmd.Flags().BoolVar(&mailboxSender, "sender", false, "access to send messages (send as & send on behalf) via SMTP")
	setaccessCmd.Flags().BoolVar(&mailboxLimitedSender, "onbehalf", false, "access to send messages (send on behalf only) via SMTP")
	setaccessCmd.Flags().Var(&notify, "notify", "notification type (all, delegates or none)")

	setaccessCmd.MarkFlagRequired("id")
	setaccessCmd.MarkFlagRequired("toId")
	setaccessCmd.MarkFlagsOneRequired("owner", "readonly", "sender", "onbehalf")
}
