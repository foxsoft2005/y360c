/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package contact

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Gets the user contact information",
	Long: `Use this command to get the user contact information.
"directory:read_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("user contact ls called")

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

		data, err := helper.GetUserById(orgId, token, userId)
		if err != nil {
			log.Fatalln("Unable to get user:", err)
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
	lsCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	lsCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	lsCmd.Flags().StringVar(&userId, "id", "", "user id")

	lsCmd.MarkFlagRequired("id")
}
