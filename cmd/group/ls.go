/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package group

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
	maxRec int
	asRaw  bool
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "gets a list of the groups",
	Long: `Use this command to get a list of Y360 groups.
"directory:read_groups" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		if !asRaw {
			log.Println("group ls called")
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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/groups?perPage=%d", helper.BaseUrl, orgId, maxRec)

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

		var data model.GroupList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		if asRaw {
			buff, _ := json.MarshalIndent(data.Groups, "", "     ")
			fmt.Print(string(buff))
		} else {
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"id", "type", "name", "email", "label", "author id", "members count"})
			for _, e := range data.Groups {
				t.AppendRow(table.Row{e.Id, e.Type, e.Name, e.Email, e.Label, e.AuthorId, e.MembersCount})
			}
			t.AppendSeparator()
			t.Render()
		}
	},
}

func init() {
	lsCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	lsCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	lsCmd.Flags().IntVar(&maxRec, "max", 100, "max records to retrieve")
	lsCmd.Flags().BoolVar(&asRaw, "raw", false, "export as raw data")
}
