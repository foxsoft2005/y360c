/*
Copyright © 2024 Kirill Chernetstky aka foxsoft2005
*/

package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/foxsoft2005/y360c/model"
	"github.com/spf13/viper"
)

const (
	BaseUrl = "https://api360.yandex.net" // Must be without trailing slash
)

// Enums "fancy" implementation

type EnumYesNo string

func (e *EnumYesNo) String() string {
	return string(*e)
}

func (e *EnumYesNo) Set(v string) error {
	switch v {
	case "yes", "no":
		*e = EnumYesNo(v)
		return nil
	default:
		return errors.New(`must be "yes" or "no"`)
	}
}

func (e *EnumYesNo) Type() string {
	return "EnumYesNo"
}

func EnumYesNoToBool(value EnumYesNo) *bool {
	var boolVar bool
	switch value {
	case "yes":
		boolVar = true
		return &boolVar
	case "no":
		boolVar = false
		return &boolVar
	default:
		return nil
	}
}

type EnumOnOff string

func (e *EnumOnOff) String() string {
	return string(*e)
}

func (e *EnumOnOff) Set(v string) error {
	switch v {
	case "on", "off":
		*e = EnumOnOff(v)
		return nil
	default:
		return errors.New(`must be "on" or "off"`)
	}
}

func (e *EnumOnOff) Type() string {
	return "EnumOnOff"
}

func EnumOnOffToBool(value EnumOnOff) *bool {
	var boolVar bool
	switch value {
	case "yes":
		boolVar = true
		return &boolVar
	case "no":
		boolVar = false
		return &boolVar
	default:
		return nil
	}
}

// Main package code

type ApiResponse struct {
	HttpCode int
	Body     []byte
}

func MakeRequest(url string, method string, token string, payload []byte) (*ApiResponse, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", token))
	if payload != nil {
		req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	}

	client := &http.Client{}
	client.Timeout = 60 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &ApiResponse{HttpCode: resp.StatusCode, Body: body}, nil
}

func GetErrorText(response *ApiResponse) error {
	if response.HttpCode == 200 {
		return nil
	}

	var errorData model.ErrorResponse
	if err := json.Unmarshal(response.Body, &errorData); err != nil {
		return fmt.Errorf("unable to parse error data (but status is %d): %s", response.HttpCode, err)
	}

	return fmt.Errorf("http %d: [%d] %s", response.HttpCode, errorData.Code, errorData.Message)
}

func GetToken() (string, error) {
	var token = viper.GetString("token")
	if token == "" {
		return "", errors.New("access token must be specified")
	}

	return token, nil
}

func GetOrgId() (int, error) {
	var orgId = viper.GetInt("org-id")
	if orgId == 0 {
		return 0, errors.New("organization id must be specified")
	}

	return orgId, nil
}

func GetMailboxById(orgId int, token string, mailboxId string) (*model.MailboxInfo, error) {
	if mailboxId == "" {
		return nil, errors.New("mailbox id must be specified")
	}

	if token == "" {
		return nil, errors.New("token must be specified")
	}

	if orgId == 0 {
		return nil, errors.New("organization id must be specified")
	}

	var url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/shared/%s", BaseUrl, orgId, mailboxId)

	resp, err := MakeRequest(url, "GET", token, nil)
	if err != nil {
		return nil, err
	}

	if err := GetErrorText(resp); err != nil {
		return nil, err
	}

	var mailbox model.MailboxInfo
	if err := json.Unmarshal(resp.Body, &mailbox); err != nil {
		return nil, err
	}

	return &mailbox, nil
}

func GetUserById(orgId int, token string, userId string) (*model.User, error) {
	if userId == "" {
		return nil, errors.New("user id must be specified")
	}

	if token == "" {
		return nil, errors.New("token must be specified")
	}

	if orgId == 0 {
		return nil, errors.New("organization id must be specified")
	}

	var url = fmt.Sprintf("%s/directory/v1/org/%d/users/%s", BaseUrl, orgId, userId)

	resp, err := MakeRequest(url, "GET", token, nil)
	if err != nil {
		return nil, err
	}

	if err := GetErrorText(resp); err != nil {
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal(resp.Body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func CheckTaskById(orgId int, token string, taskId string) (*model.TaskStatusResponse, error) {
	if taskId == "" {
		return nil, errors.New("task id must be specified")
	}

	if token == "" {
		return nil, errors.New("token must be specified")
	}

	if orgId == 0 {
		return nil, errors.New("organization id must be specified")
	}

	var url = fmt.Sprintf("%s/admin/v1/org/%d/mailboxes/tasks/%s", BaseUrl, orgId, taskId)

	resp, err := MakeRequest(url, "GET", token, nil)
	if err != nil {
		return nil, err
	}

	if err := GetErrorText(resp); err != nil {
		return nil, err
	}

	var task model.TaskStatusResponse
	if err := json.Unmarshal(resp.Body, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func Confirm(message string) bool {

	var input string

	log.Print(message)

	_, err := fmt.Scanln(&input)
	if err != nil {
		log.Fatalln("Failed to parse keyboard input:", err)
	}

	input = strings.ToLower(input)

	if input == "y" || input == "yes" {
		return true
	}

	return false
}

func ToNullableString(value string) *string {
	var empty = ""

	if value == "" {
		return nil
	}

	if value == "<clear>" {
		return &empty
	}

	return &value
}
