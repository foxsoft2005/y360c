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

type mailboxType string

const sharedMailbox mailboxType = "shared"

func (e *mailboxType) String() string {
	return string(*e)
}

func (e *mailboxType) Set(v string) error {
	switch v {
	case "shared", "delegated":
		*e = mailboxType(v)
		return nil
	default:
		return errors.New(`must be one of "shared", or "delegated"`)
	}
}

func (e *mailboxType) Type() string {
	return "mailboxType"
}

var (
	maxRec      int
	requestType mailboxType
	asRaw       bool
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Gets a list of the shared or delegated mailboxes",
	Long: `Use this command to retrieve a list of shared or delegated mailboxes.
"ya360_admin:mail_read_shared_mailbox_inventory" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		if !asRaw {
			log.Print("mailbox ls called")
		}

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

		var url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/%s?perPage=%d", helper.BaseUrl, orgId, requestType, maxRec)

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

		var data model.MailboxList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		if asRaw {
			buff, _ := json.MarshalIndent(data.Items, "", "     ")
			fmt.Print(string(buff))
		} else {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"resource", "count"})
			for _, e := range data.Items {
				t.AppendRow(table.Row{e.ResourceId, e.Count})
			}
			t.AppendSeparator()
			t.Render()
		}

	},
}

func init() {
	requestType = sharedMailbox

	lsCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	lsCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	lsCmd.Flags().IntVar(&maxRec, "max", 10, "max records to retrieve")
	lsCmd.Flags().Var(&requestType, "type", "type of the mailbox (shared or delegated)")
	lsCmd.Flags().BoolVar(&asRaw, "raw", false, "export as raw data")
}
