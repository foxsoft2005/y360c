/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package dept

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

type sortingMethod string

const byId sortingMethod = "id"

func (e *sortingMethod) String() string {
	return string(*e)
}

func (e *sortingMethod) Set(v string) error {
	switch v {
	case "id", "name":
		*e = sortingMethod(v)
		return nil
	default:
		return errors.New(`must be one of "id", or "name"`)
	}
}

func (e *sortingMethod) Type() string {
	return "sortingMethod"
}

// Command flags & parameters
var maxRec int
var orderBy sortingMethod

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "gets a list of the departments",
	Long: `Use this command to retrieve a list of departments of selected organization.
"directory:read_departments" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("dept ls called")

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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/departments?perPage=%d&orderBy=%s", helper.BaseUrl, orgId, maxRec, orderBy)
		if parentId > 0 {
			url += fmt.Sprintf("&parentId=%d", parentId)
		}

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

		var data model.DepartmentList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"id", "parent id", "name", "description", "email", "label", "members count"})
		for _, e := range data.Departments {
			t.AppendRow(table.Row{e.Id, e.ParentId, e.Name, e.Description, e.Email, e.Label, e.MembersCount})
		}

		t.AppendSeparator()
		t.Render()

	},
}

func init() {
	orderBy = byId

	lsCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	lsCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	lsCmd.Flags().IntVar(&maxRec, "max", 100, "max records to retrieve")
	lsCmd.Flags().IntVar(&parentId, "parentId", 0, "parent depratment id")
	lsCmd.Flags().Var(&orderBy, "orderBy", "sort by (id or name)")
}
