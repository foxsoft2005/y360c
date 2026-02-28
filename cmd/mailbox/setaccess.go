// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package mailbox

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/goccy/go-json"

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
	toEmail              string
	toId                 string
	notify               notificationType
	mailboxOwner         bool
	mailboxImap          bool
	mailboxSender        bool
	mailboxLimitedSender bool
	clearAccess          bool
	checkStatus          bool
)

// accessCmd represents the delegate command
var accessCmd = &cobra.Command{
	Use:   "set-access",
	Short: "Set access to the mailbox to other user",
	Long: `Use this command to set access to the mailbox to other user.
"ya360_admin:mail_write_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("mailbox set-access called")

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

		if mailboxName != "" {
			us, err := helper.GetUserByEmail(orgId, token, mailboxName)
			if err != nil {
				log.Fatalln("Failed to get user by email", err)
			}

			if us == nil {
				log.Fatalf("User (mailbox) %s does not found", mailboxName)
			}

			mailboxId = us.Id
		}

		if toEmail != "" {
			us, err := helper.GetUserByEmail(orgId, token, toEmail)
			if err != nil {
				log.Fatalln("Failed to get user by email", err)
			}

			if us == nil {
				log.Fatalf("User (mailbox) %s does not found", toEmail)
			}

			toId = us.Id
		}

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/set/%s?actorId=%s&notify=%s", helper.BaseUrl, orgId, mailboxId, toId, notify)

		var method = "POST"
		var entry model.MailAccessSettings

		if !clearAccess {
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
		}

		payload, _ := json.Marshal(entry)
		resp, err := helper.MakeRequest(url, method, token, payload)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
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

		if checkStatus {
			log.Printf("Checking status for task %s", data.TaskId)

			time.Sleep(2 * time.Second)
			task, err := helper.CheckTaskById(orgId, token, data.TaskId)
			if err != nil {
				log.Fatalln("Unable to check task:", err)
			}

			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendRow(table.Row{"Task Id", data.TaskId})
			t.AppendRow(table.Row{"Status", task.Status})
			t.AppendSeparator()
			t.Style().Options.SeparateRows = true
			t.Render()
		}
	},
}

func init() {
	notify = notifyAll

	accessCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	accessCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	accessCmd.Flags().StringVar(&mailboxId, "id", "", "mailbox (or user) id")
	accessCmd.Flags().StringVar(&mailboxName, "email", "", "mailbox (or user) email address")
	accessCmd.Flags().StringVar(&toId, "to-id", "", "user id to whom access is delegated")
	accessCmd.Flags().StringVar(&toEmail, "to-email", "", "user email address to whom access is delegated")
	accessCmd.Flags().BoolVar(&mailboxOwner, "owner", false, "full access to the mailbox (includes all other access types)")
	accessCmd.Flags().BoolVar(&mailboxImap, "reader", false, "access to read messages via IMAP")
	accessCmd.Flags().BoolVar(&mailboxSender, "sender", false, "access to send messages (send as & send on behalf) via SMTP")
	accessCmd.Flags().BoolVar(&mailboxLimitedSender, "on-behalf", false, "access to send messages (send on behalf only) via SMTP")
	accessCmd.Flags().Var(&notify, "notify", "notification type (all, delegates or none)")
	accessCmd.Flags().BoolVar(&checkStatus, "check-status", false, "automatically check task status if possible")
	accessCmd.Flags().BoolVar(&clearAccess, "clearAccess", false, "clearAccess access options for the selected user if any")

	accessCmd.MarkFlagsOneRequired("id", "email")
	accessCmd.MarkFlagsOneRequired("to-id", "to-email")

	accessCmd.MarkFlagsOneRequired("owner", "reader", "sender", "on-behalf", "clearAccess")

	accessCmd.MarkFlagsMutuallyExclusive("id", "email")
	accessCmd.MarkFlagsMutuallyExclusive("to-id", "to-email")

	accessCmd.MarkFlagsMutuallyExclusive("owner", "reader")
	accessCmd.MarkFlagsMutuallyExclusive("owner", "sender")
	accessCmd.MarkFlagsMutuallyExclusive("owner", "on-behalf")

	accessCmd.MarkFlagsMutuallyExclusive("clearAccess", "owner")
	accessCmd.MarkFlagsMutuallyExclusive("clearAccess", "sender")
	accessCmd.MarkFlagsMutuallyExclusive("clearAccess", "reader")
	accessCmd.MarkFlagsMutuallyExclusive("clearAccess", "on-behalf")
}
