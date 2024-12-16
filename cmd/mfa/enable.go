/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/
package mfa

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

type validationMethod string

const (
	smsMethod validationMethod = "sms"
)

func (e *validationMethod) String() string {
	return string(*e)
}

func (e *validationMethod) Set(v string) error {
	switch v {
	case "sms", "phone":
		*e = validationMethod(v)
		return nil
	default:
		return errors.New(`must be one of "sms", or "phone"`)
	}
}

func (e *validationMethod) Type() string {
	return "validationMethod"
}

var duration int
var logoutUsers bool
var validation validationMethod

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enables mandatory 2FA for domain users",
	Long: `Use this command to enable madatory two-factor authentication (2FA) for the selected organization.
"ya360_security:domain_2fa_write" permission is required (see Y360 help topics).`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("mfa enable called")

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

		var url = fmt.Sprintf("%s/security/v1/org/%d/domain_2fa", helper.BaseUrl, orgId)

		payload, _ := json.Marshal(model.MfaActivation{Duration: duration, LogoutUsers: logoutUsers, ValidationMethod: validation.String()})

		resp, err := helper.MakeRequest(url, "POST", token, payload)
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

		var data model.MfaSetup
		if err := json.Unmarshal(resp.Body, &data); err != nil {
			log.Fatalln("Unable to evaluate data:", err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendRow(table.Row{"Duration", data.Duration})
		t.AppendRow(table.Row{"Enabled", data.Enabled})
		t.AppendRow(table.Row{"Enabled At", data.EnabledAt})
		t.AppendSeparator()
		t.Style().Options.SeparateRows = true
		t.Render()
	},
}

func init() {
	validation = smsMethod

	enableCmd.Flags().IntVarP(&orgId, "orgId", "o", 0, "organization id")
	enableCmd.Flags().StringVarP(&token, "token", "t", "", "access token")
	enableCmd.Flags().IntVar(&duration, "duration", 0, "period (sec.) to postpone MFA setup by user")
	enableCmd.Flags().BoolVar(&logoutUsers, "logout", false, "logout all users on MFA activation")
	enableCmd.Flags().Var(&validation, "validate", "validation method (sms, phone)")

	enableCmd.MarkFlagRequired("duration")
}
