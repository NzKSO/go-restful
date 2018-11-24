package restful

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client ...
type Client struct {
	client     *http.Client
	header     http.Header
	query      url.Values
	body       io.Reader
	baseURL    string
	requestURI string
	accessKey  string
	secretKey  string
}

// NewClient ...
func NewClient() *Client {
	return &Client{
		client: &http.Client{},
		header: http.Header{},
		query:  url.Values{},
	}
}

// SetBaseURL ...
func (c *Client) SetBaseURL(baseURL string) *Client {
	c.baseURL = baseURL
	return c
}

// SetHeader ...
func (c *Client) SetHeader(key, value string) *Client {
	c.header.Set(key, value)
	return c
}

// SetHeaders ...
func (c *Client) SetHeaders(header map[string]string) *Client {
	for k, v := range header {
		c.header.Set(k, v)
	}
	return c
}

// SetProxy ...
func (c *Client) SetProxy(proxy string) {
	p, err := url.ParseRequestURI(proxy)
	if err == nil {
		switch p.Scheme {
		case "http", "https", "socks5":
		default:
			return
		}

		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = http.ProxyURL(p)
		c.client.Transport = transport
	}
}

// SetBody ...
func (c *Client) SetBody(body io.Reader) {
	c.body = body
}

// SetQuery ...
func (c *Client) SetQuery(key, value string) *Client {
	c.query.Set(key, value)
	return c
}

// SetQuerys ...
func (c *Client) SetQuerys(params map[string]string) *Client {
	for k, v := range params {
		c.query.Set(k, v)
	}
	return c
}

// SetRequestURI ...
func (c *Client) SetRequestURI(path string) {
	c.requestURI = path
}

// Do ...
func (c *Client) Do(method string) ([]byte, error) {
	requestURL := c.baseURL + c.requestURI
	if c.query != nil && len(c.query) > 0 {
		requestURL += "?" + c.query.Encode()
	}
	req, err := http.NewRequest(method, requestURL, c.body)
	if err != nil {
		return nil, err
	}
	req.Header = c.header
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
