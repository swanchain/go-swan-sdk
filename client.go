package swan

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

type APIClient struct {
	apiKey     string
	httpClient *HttpClient
}

func NewAPIClient(apiKey string, testnet ...bool) *APIClient {
	host := gatewayMainnet
	if len(testnet) > 0 && testnet[0] {
		host = gatewayTestnet
	}

	header := make(http.Header)
	header.Add("Authorization", "Bearer "+apiKey)

	return &APIClient{
		apiKey:     apiKey,
		httpClient: NewHttpClient(host, header),
	}
}

func (c *APIClient) Hardwares() ([]*Hardware, error) {
	var result HardwareResult

	if err := c.httpClient.Get(apiMachines, nil, NewResult(&result)); err != nil {
		return nil, err
	}

	return result.Hardware, nil
}

func (c *APIClient) TaskInfo(taskUUID string) (*TaskInfo, error) {
	var result TaskInfo

	if err := c.httpClient.Get(fmt.Sprintf("%s/%s", apiTask, taskUUID), nil, NewResult(&result)); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *APIClient) Tasks(req *TaskQueryReq) (total int64, list []*TaskInfo, err error) {
	api := apiTasks
	if req != nil {
		api += fmt.Sprintf("?wallet=%s&size=%d&page=%d", req.Wallet, req.Size, req.Page)

	}

	var result PageResult
	result.List = &list

	if err = c.httpClient.Get(api, nil, NewResult(&result)); err != nil {
		return
	}
	total = result.Total
	return
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
