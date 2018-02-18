package wxweb

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Client struct {
	client    *http.Client
	userAgent string
}

type Header map[string]string

func NewClient() *Client {
	var netTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial:  (&net.Dialer{Timeout: 100 * time.Second}).Dial,
		// TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		// TLSHandshakeTimeout: 100 * time.Second,
	}
	cookieJar, _ := cookiejar.New(nil)

	httpClient := &http.Client{
		Timeout:   time.Second * 100,
		Transport: netTransport,
		Jar:       cookieJar,
	}

	return &Client{
		client:    httpClient,
		userAgent: "ApiV2 Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36 ",
	}
}
func (c *Client) SetJar(jar http.CookieJar) {
	c.client.Jar = jar
}
func (c *Client) Get(url string, data *url.Values) ([]byte, error) {
	if data != nil {
		url = url + "?" + data.Encode()
	}
	return c.fetch("GET", url, []byte(""), Header{})
}
func (c *Client) GetByte(url string, data []byte) ([]byte, error) {

	return c.fetch("GET", url, data, Header{})
}
func (c *Client) GetWithHeader(url string, heder Header) ([]byte, error) {

	return c.fetch("GET", url, nil, heder)
}

func (c *Client) Post(url string, data *url.Values) ([]byte, error) {
	return c.fetch("POST", url, []byte(data.Encode()), Header{"Content-Type": "application/x-www-form-urlencoded"})
}

func (c *Client) PostJson(url string, m map[string]interface{}) ([]byte, error) {
	jsonString, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return c.fetch("POST", url, jsonString, Header{"Content-Type": "application/json; charset=UTF-8"})
}
func (c *Client) PostJsonByte(url string, json []byte) ([]byte, error) {

	return c.fetch("POST", url, json, Header{"Content-Type": "application/json; charset=UTF-8"})
}
func (c *Client) PostJsonByteForResp(url string, json []byte) (*http.Response, []byte, error) {
	return c.fetchResp("POST", url, json, Header{"Content-Type": "application/json; charset=UTF-8"})
}
func (c *Client) fetchReponse(method string, uri string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.client.Do(req)
}

func (c *Client) fetchReponseWithReader(method string, uri string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.client.Do(req)
}
func (c *Client) fetchWithReader(method string, uri string, body io.Reader, headers Header) ([]byte, error) {
	resp, err := c.fetchReponseWithReader(method, uri, body, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Client) fetch(method string, uri string, body []byte, headers Header) ([]byte, error) {
	resp, err := c.fetchReponse(method, uri, body, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Client) fetchResp(method string, uri string, body []byte, headers Header) (resp *http.Response, b []byte, err error) {
	resp, err = c.fetchReponse(method, uri, body, headers)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	b, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return nil, nil, err2
	}
	return resp, b, nil
}
