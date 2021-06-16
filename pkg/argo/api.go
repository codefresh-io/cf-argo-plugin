package argo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

func buildHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

type argo struct {
	Host     string
	Username string
	Password string
	Token    string
}

type Argo interface {
	GetLatestHistoryId(application string) (int64, error)
}

type ClientOptions struct {
	Host     string
	Username string
	Password string
	Token    string
}

func New(options *ClientOptions) Argo {
	return &argo{
		Host:     options.Host,
		Username: options.Username,
		Password: options.Password,
		Token:    options.Token,
	}
}

type requestOptions struct {
	path   string
	method string
}

type History struct {
	Status struct {
		History []struct {
			Id int64 `json:"id"`
		} `json:"history"`
	} `json:"status"`
}

func (c *argo) GetLatestHistoryId(application string) (int64, error) {

	options := &requestOptions{
		path:   "/api/v1/applications/" + application,
		method: "GET",
	}

	result := &History{}
	_ = c.requestAPI(options, result)

	historyList := result.Status.History
	if len(historyList) == 0 {
		return -1, nil
	}

	return historyList[len(historyList)-1].Id, nil
}

func (c *argo) getToken() string {

	client := buildHttpClient()

	message := map[string]interface{}{
		"username": c.Username,
		"password": c.Password,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.Post(c.Host+"/api/v1/session", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Println(err)
	}
	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	defer resp.Body.Close()

	return result["token"].(string)
}

func (c *argo) requestAPI(opt *requestOptions, target interface{}) error {

	token := c.Token
	if token == "" {
		token = c.getToken()
	}

	client := buildHttpClient()

	var body []byte
	finalURL := fmt.Sprintf("%s%s", c.Host, opt.path)

	request, err := http.NewRequest(opt.method, finalURL, bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		argoError := make(map[string]interface{})
		err = json.NewDecoder(response.Body).Decode(argoError)

		if err != nil {
			return err
		}

		return fmt.Errorf("%d: %v", response.StatusCode, argoError)
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(target)

	if err != nil {
		return err
	}

	return nil
}
