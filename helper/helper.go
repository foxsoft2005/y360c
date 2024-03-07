/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/

package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/foxsoft2005/y360c/model"
	"github.com/spf13/viper"
)

const (
	BaseUrl = "https://api360.yandex.net" // Must be without trailing slash
)

type ApiResponse struct {
	HttpCode int
	Body     []byte
}

func MakeRequest(url string, method string, token string, payload []byte) (*ApiResponse, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
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

func GetToken() (string, error) {
	var token = viper.GetString("token")
	if token == "" {
		return "", errors.New("access token must be specified")
	}

	return token, nil
}

func GetOrgId() (int, error) {
	var orgId = viper.GetInt("orgId")
	if orgId == 0 {
		return 0, errors.New("organization id must be specified")
	}

	return orgId, nil
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

	if resp.HttpCode != 200 {
		var error model.ErrorResponse
		if err := json.Unmarshal(resp.Body, &error); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("http %d: [%d] %s", resp.HttpCode, error.Code, error.Message)
	}

	var user model.User
	if err := json.Unmarshal(resp.Body, &user); err != nil {
		return nil, err
	}

	return &user, nil
}
