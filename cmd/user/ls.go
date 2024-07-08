/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package user

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/foxsoft2005/y360c/helper"
	"github.com/foxsoft2005/y360c/model"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var (
	userName string
	email    string
	deptId   int
	maxRec   int
	asCsv    bool
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "gets a list of the users",
	Long: `Use this command to retrieve a list of the users of selected organization.
"directory:read_users" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		if !asRaw {
			log.Print("user ls called")
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

		var url = fmt.Sprintf("%s/directory/v1/org/%d/users?perPage=%d", helper.BaseUrl, orgId, maxRec)

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

		var data model.UserList
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		var users []model.User
		showAll := (userId == "") && (email == "") && (userName == "") && (deptId == 0)
		found := false

		r, _ := regexp.Compile(fmt.Sprintf("(?i)%s", userName))

		for _, e := range data.Users {
			switch true {
			case userId != "":
				found = strings.EqualFold(e.Id, userId)
			case email != "":
				found = strings.EqualFold(e.Email, email)
			case userName != "":
				found = r.MatchString(e.Name.First) || r.MatchString(e.Name.Last)
			case deptId != 0:
				found = deptId == e.DepartmentId
			}
			if found || showAll {
				users = append(users, e)
			}
		}

		if asRaw {
			buff, _ := json.MarshalIndent(users, "", "     ")
			fmt.Print(string(buff))
		} else {
			t := table.NewWriter()
			if !asCsv {
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"id", "name", "email", "timezone", "enabled", "dismissed", "admin"})
			} else {
				t.AppendHeader(
					table.Row{
						"Id",
						"External Id",
						"Name",
						"Nickname",
						"Display Name",
						"About",
						"Gender",
						"Birthday",
						"Email",
						"Avatar Id",
						"Position",
						"Department",
						"Timezone",
						"Aliases",
						"Groups",
						"Enabled",
						"Dismissed",
						"Admin",
						"Robot",
						"Created At",
						"Updated At",
					},
				)
			}

			for _, e := range users {
				if !asCsv {
					t.AppendRow(
						table.Row{e.Id, strings.TrimSpace(fmt.Sprintf("%s %s %s", e.Name.Last, e.Name.First, e.Name.Middle)), e.Email, e.Timezone, e.IsEnabled, e.IsDismissed, e.IsAdmin},
					)
				} else {
					t.AppendRow(
						table.Row{
							e.Id,
							e.ExternalId,
							strings.TrimSpace(fmt.Sprintf("%s %s %s", e.Name.Last, e.Name.First, e.Name.Middle)),
							e.Nickname,
							e.DisplayName,
							e.About,
							e.Gender,
							e.Birthday,
							e.Email,
							e.AvararId,
							e.Position,
							e.DepartmentId,
							e.Timezone,
							strings.Join(e.Aliases, ","),
							strings.Trim(strings.Join(strings.Fields(fmt.Sprint(e.Groups)), ","), "[]"),
							e.IsEnabled,
							e.IsDismissed,
							e.IsAdmin,
							e.IsRobot,
							e.CreatedAt,
							e.UpdatedAt,
						},
					)
				}
			}
			t.AppendSeparator()
			if asCsv {
				path, _ := os.Getwd()
				fileName := filepath.Join(path, fmt.Sprintf("users_%d.csv", orgId))

				f, err := os.Create(fileName)
				if err != nil {
					log.Fatalln("Unable to create csv file:", err)
				}
				defer f.Close()

				_, err1 := f.WriteString(t.RenderCSV())
				if err1 != nil {
					log.Fatalln("Unable to write to csv file:", err1)
				}

				log.Printf("The users where successfully exported to %s", fileName)
			} else {
				t.Render()
			}
		}
	},
}

func init() {
	listCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	listCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	listCmd.Flags().IntVar(&maxRec, "max", 1000, "max records to retrieve")
	listCmd.Flags().StringVar(&userId, "id", "", "find user by id")
	listCmd.Flags().StringVar(&email, "email", "", "find user by email")
	listCmd.Flags().StringVar(&userName, "name", "", "find user(s) by name")
	listCmd.Flags().BoolVar(&asCsv, "csv", false, "export data to csv-file")
	listCmd.Flags().IntVar(&deptId, "deptId", 0, "find user(s) by department id")
	listCmd.Flags().BoolVar(&asRaw, "raw", false, "export as raw data")

	listCmd.MarkFlagsMutuallyExclusive("id", "email", "name", "deptId")
	listCmd.MarkFlagsMutuallyExclusive("csv", "raw")

}
