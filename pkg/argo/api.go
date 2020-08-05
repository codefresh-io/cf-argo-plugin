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

type Argo struct {
	Host     string
	Username string
	Password string
}

type requestOptions struct {
	path   string
	method string
}

type History struct {
	Status struct {
		History []struct {
			Revision string `json:"revision"`
		} `json:"history"`
	} `json:"status"`
}

func (c *Argo) GetLatestHistoryRevision(application string) (string, error) {

	options := &requestOptions{
		path:   "/api/v1/applications/" + application,
		method: "GET",
	}

	result := &History{}
	_ = c.requestAPI(options, result)

	historyList := result.Status.History
	return historyList[len(historyList)-1].Revision, nil
}

func (c *Argo) getToken() string {

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

func (c *Argo) requestAPI(opt *requestOptions, target interface{}) error {
	token := c.getToken()

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
