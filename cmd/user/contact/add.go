// Copyright © 2024-2026 Kirill Chernetsky aka foxsoft2005

package contact

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-json"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

type contactInfoType string

func (e *contactInfoType) String() string {
	return string(*e)
}

func (e *contactInfoType) Set(v string) error {
	switch v {
	case "email", "phone_extension", "phone", "site", "icq", "twitter", "skype":
		*e = contactInfoType(v)
		return nil
	default:
		return errors.New(`must be one of "email", "phone_extension", "phone", "site", "icq", "twitter", "skype"`)
	}
}

func (e *contactInfoType) Type() string {
	return "contactInfoType"
}

var (
	contactType  contactInfoType
	contactLabel string
	contactValue string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a contact information",
	Long: `Use this command to add a contact information for the selected user.
"directory:write_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user contact add called")

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

		if userEmail != "" {
			us, err := helper.GetUserByEmail(orgId, token, userEmail)
			if err != nil {
				log.Fatalln("Failed to get user by email", err)
			}

			if us == nil {
				log.Fatalf("User (mailbox) %s does not found", userEmail)
			}

			userId = us.Id
		}

		var url = fmt.Sprintf("%s/directory/v1/org/%d/users/%s/contacts", helper.BaseUrl, orgId, userId)

		user, err := helper.GetUserById(orgId, token, userId)
		if err != nil {
			log.Fatalln("Unable to get user:", err)
		}

		// getting all existing items
		var entry model.ContactInfoList
		for _, item := range user.Contacts {
			if !item.Synthetic {
				entry.Items = append(entry.Items, model.ContactInfo{Type: item.Type, Value: item.Value, Label: item.Label})
			}
		}

		// adding new
		entry.Items = append(entry.Items, model.ContactInfo{Type: string(contactType), Value: contactValue, Label: contactLabel})

		payload, _ := json.Marshal(entry)
		resp, err := helper.MakeRequest(url, "PUT", token, payload)
		if err != nil {
			log.Fatalln("Unable to make API request:", err)
		}

		if err := helper.GetErrorText(resp); err != nil {
			log.Fatalln(err)
		}

		var data model.User
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Id", data.Id, ""})
		t.AppendRow(table.Row{"Name", strings.TrimSpace(fmt.Sprintf("%s %s %s", data.Name.Last, data.Name.First, data.Name.Middle)), ""})
		for index, item := range data.Contacts {
			if index == 0 {
				t.AppendRow(table.Row{"Contact info", "", ""})
			}

			var s = item.Type
			if item.Main {
				s = fmt.Sprintf("%s (main)", item.Type)
			}

			var s1 = item.Value
			if item.Synthetic {
				s1 = fmt.Sprintf("%s (readonly)", item.Value)
			}

			t.AppendRow(table.Row{"", s, s1})
		}
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	addCmd.Flags().IntVarP(&orgId, "org-id", "o", 0, "organization id")
	addCmd.Flags().StringVarP(&token, "token", "t", "", "access token")

	addCmd.Flags().StringVar(&userId, "id", "", "user id")
	addCmd.Flags().StringVar(&userEmail, "email", "", "user email address")

	addCmd.Flags().Var(&contactType, "type", "entry type (email, phone, phone_extension, etc.)")
	addCmd.Flags().StringVar(&contactValue, "value", "", "entry value")
	addCmd.Flags().StringVar(&contactLabel, "label", "", "entry label")

	addCmd.MarkFlagsOneRequired("id", "email")
	addCmd.MarkFlagsMutuallyExclusive("id", "email")

	err := addCmd.MarkFlagRequired("type")
	if err != nil {
		log.Fatalln("Error marking flag as required:", err)
	}
	err = addCmd.MarkFlagRequired("value")
	if err != nil {
		log.Fatalln("Error marking flag as required:", err)
	}
}
