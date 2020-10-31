package codefresh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Codefresh interface {
	GetIntegration(name string) (*ArgoIntegration, error)
	StartSyncTask(name string) (*TaskResult, error)
	SendMetadata(metadata *ArgoApplicationMetadata) error
	RollbackToStable(name string) (*TaskResult, error)
	requestAPI(*requestOptions, interface{}) error
}

type CodefreshError struct {
	Status  int         `json:"status"`
	Code    string      `json:"code"`
	Name    string      `json:"name"`
	Message string      `json:"message"`
	Context interface{} `json:"context"`
}

type ClientOptions struct {
	Token string
	Host  string
}

type ArgoIntegration struct {
	Type string              `json:"type"`
	Data ArgoIntegrationData `json:"data"`
}

type ArgoApplicationMetadata struct {
	Pipeline        string `json:"pipeline"`
	BuildId         string `json:"buildId"`
	HistoryId       int64  `json:"historyId"`
	ApplicationName string `json:"name"`
}

type TaskResult struct {
	BuildId string `json:"id"`
}

type ArgoIntegrationData struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type requestOptions struct {
	path   string
	method string
	body   []byte
}

type codefresh struct {
	host   string
	token  string
	client *http.Client
}

func New(opt *ClientOptions) Codefresh {
	return &codefresh{
		host:   opt.Host,
		token:  opt.Token,
		client: &http.Client{},
	}
}

func (c *codefresh) GetIntegration(name string) (*ArgoIntegration, error) {
	r := &ArgoIntegration{}
	err := c.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/argo/%s", name),
		method: "GET",
	}, r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *codefresh) SendMetadata(metadata *ArgoApplicationMetadata) error {
	metadataBytes := new(bytes.Buffer)
	json.NewEncoder(metadataBytes).Encode(metadata)

	var result map[string]interface{}

	err := c.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/environments-v2/argo/metadata"),
		method: "POST",
		body:   metadataBytes.Bytes(),
	}, &result)

	if err != nil {
		return err
	}

	return nil
}

func (c *codefresh) StartSyncTask(name string) (*TaskResult, error) {
	r := &TaskResult{}
	err := c.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/environments-v2/sync/%s", name),
		method: "GET",
	}, r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *codefresh) RollbackToStable(name string) (*TaskResult, error) {
	r := &TaskResult{}
	err := c.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/gitops/argocd/%s/rollbackToStable", name),
		method: "GET",
	}, r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *codefresh) requestAPI(opt *requestOptions, target interface{}) error {
	finalURL := fmt.Sprintf("%s%s", c.host, opt.path)

	request, err := http.NewRequest(opt.method, finalURL, bytes.NewBuffer(opt.body))

	if err != nil {
		return err
	}

	request.Header.Set("Authorization", c.token)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		cfError := &CodefreshError{}
		err = json.NewDecoder(response.Body).Decode(cfError)

		if err != nil {
			return err
		}

		return fmt.Errorf("%d: %s", response.StatusCode, cfError.Message)
	}

	err = json.NewDecoder(response.Body).Decode(target)

	if err != nil {
		return err
	}

	return nil
}
