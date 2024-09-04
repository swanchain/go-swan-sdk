package swan

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type HttpClient struct {
	host   string
	header http.Header
	client *http.Client
}

func NewHttpClient(host string, header http.Header, client ...*http.Client) *HttpClient {
	hc := httpClient
	if len(client) > 0 {
		hc = client[0]
	}
	return &HttpClient{
		host:   host,
		header: header,
		client: hc,
	}
}

func (c *HttpClient) Get(api string, queries url.Values, dest any) error {
	if queries != nil {
		api += "?" + queries.Encode()
	}
	return c.Request(http.MethodGet, api, nil, dest, "")
}

func (c *HttpClient) PostForm(api string, data url.Values, dest any) error {
	return c.Request(http.MethodPost, api, strings.NewReader(data.Encode()), dest, "application/x-www-form-urlencoded")
}

func (c *HttpClient) PostJSON(api string, data any, dest any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.Request(http.MethodPost, api, bytes.NewReader(b), dest, "application/json")
}

func (c *HttpClient) Request(method string, api string, body io.Reader, dest any, contentType ...string) (err error) {
	//paras := ""
	if body != nil {
		rb, _ := io.ReadAll(body)
		//paras = string(rb)
		body = bytes.NewReader(rb)
	}

	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return errors.New("dest is not a pointer")
	}

	url := c.host
	if strings.HasPrefix(api, "/") {
		url += api
	} else {
		url += "/" + api
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	for key := range c.header {
		req.Header.Set(key, c.header.Get(key))
	}
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType[0])
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	//fmt.Printf("method: %s, api: %s, paras: %s, response: %s\n", method, api, paras, string(bd))

	if err = json.Unmarshal(bd, dest); err != nil {
		return
	}

	if checker, ok := dest.(Checker); ok {
		return checker.Check()
	}
	return
}

var httpClient = &http.Client{
	Timeout: time.Second * 30,
}

type Checker interface {
	Check() error
}
