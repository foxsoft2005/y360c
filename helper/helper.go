/*
Copyright © 2024 Kirill Chernetsky <foxsoft2005@gmail.com>
*/

package helper

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

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
