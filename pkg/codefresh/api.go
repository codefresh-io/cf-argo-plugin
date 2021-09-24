package codefresh

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type Codefresh interface {
	GetIntegration(name string) (*ArgoIntegration, error)
	StartSyncTask(name string) (*TaskResult, error)
	SendMetadata(metadata *ArgoApplicationMetadata) (error, []UpdatedActivity)
	RollbackToStable(name string, payload Rollback) (*TaskResult, error)
	GetEnvironments() ([]CFEnvironment, error)
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

type MongoCFEnvWrapper struct {
	Docs []CFEnvironment `json:"docs"`
}

type CFEnvironment struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Type        string `json:"type"`
		Application string `json:"application"`
	} `json:"spec"`
}

type ArgoApplicationMetadata struct {
	Pipeline        string `json:"pipeline"`
	BuildId         string `json:"buildId"`
	HistoryId       int64  `json:"historyId"`
	ApplicationName string `json:"name"`
}

type Rollback struct {
	ContextName     string `json:"contextName"`
	ApplicationName string `json:"applicationName"`
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

type updateMetadataResponse struct {
	Activities []UpdatedActivity `json:"activities"`
}

type UpdatedActivity struct {
	ActivityId      string `json:"_id"`
	EnvironmentId   string `json:"environmentId"`
	EnvironmentName string `json:"environmentName"`
	ApplicationName string `json:"applicationName"`
}

type codefresh struct {
	host   string
	token  string
	client *http.Client
}

func New(opt *ClientOptions) Codefresh {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &codefresh{
		host:   opt.Host,
		token:  opt.Token,
		client: &http.Client{Transport: tr},
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

func (c *codefresh) SendMetadata(metadata *ArgoApplicationMetadata) (error, []UpdatedActivity) {
	metadataBytes := new(bytes.Buffer)
	json.NewEncoder(metadataBytes).Encode(metadata)

	var result updateMetadataResponse

	err := c.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/environments-v2/argo/metadata"),
		method: "POST",
		body:   metadataBytes.Bytes(),
	}, &result)

	if err != nil {
		return err, nil
	}
	return nil, result.Activities
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

func (c *codefresh) RollbackToStable(name string, payload Rollback) (*TaskResult, error) {
	metadataBytes := new(bytes.Buffer)
	json.NewEncoder(metadataBytes).Encode(payload)

	r := &TaskResult{}
	err := c.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/gitops/argocd/%s/rollbackToStable", name),
		method: "POST",
		body:   metadataBytes.Bytes(),
	}, r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *codefresh) GetEnvironments() ([]CFEnvironment, error) {
	var result MongoCFEnvWrapper
	err := c.requestAPI(&requestOptions{
		method: "GET",
		path:   "/api/gitops/application?plain=true&isEnvironment=false",
	}, &result)
	if err != nil {
		return nil, err
	}

	return result.Docs, nil
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
