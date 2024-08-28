package swan

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
)

type APIClient struct {
	apiKey     string
	httpClient *HttpClient
}

func NewAPIClient(apiKey string, testnet ...bool) (*APIClient, error) {
	host := gatewayMainnet
	if len(testnet) > 0 && testnet[0] {
		host = gatewayTestnet
	}

	client := &APIClient{
		apiKey: apiKey,
		httpClient: &HttpClient{
			host: host,
		},
	}

	token, err := client.login()
	if err != nil {
		return nil, err
	}

	header := make(http.Header)
	header.Add("Authorization", "Bearer "+token)
	client.httpClient.header = header
	return client, nil
}

func (c *APIClient) login() (string, error) {
	var token string

	if err := c.httpClient.PostForm(apiLogin, url.Values{"api_key": {c.apiKey}}, NewResult(&token)); err != nil {
		return "", err
	}

	return token, nil
}

func (c *APIClient) Hardwares() ([]*Hardware, error) {
	var result HardwareResult

	if err := c.httpClient.Get(apiMachines, nil, NewResult(&result)); err != nil {
		return nil, err
	}

	return result.Hardware, nil
}

func (c *APIClient) TaskInfo(taskUUID string) (*TaskDetails, error) {
	var result TaskDetails

	if err := c.httpClient.Get(fmt.Sprintf("%s/%s", apiTask, taskUUID), nil, NewResult(&result)); err != nil {
		return nil, err
	}

	return &result, nil
}

type Result struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func NewResult(dest any) *Result {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		panic("dest must be a pointer")
	}
	var result Result
	result.Data = dest
	return &result
}

func (r *Result) Check() error {
	if r.Status != "success" {
		return errors.New(r.Message)
	}
	return nil
}
