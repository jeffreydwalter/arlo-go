package request

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/pkg/errors"
)

type Client struct {
	BaseURL        *url.URL
	BaseHttpHeader http.Header
	httpClient     http.Client
}

func NewClient(baseurl string) (*Client, error) {
	var err error
	var jar *cookiejar.Jar

	options := cookiejar.Options{}

	if jar, err = cookiejar.New(&options); err != nil {
		return nil, errors.Wrap(err, "failed to create client object")
	}

	var u *url.URL
	if u, err = url.Parse(baseurl); err != nil {
		return nil, errors.Wrap(err, "failed to create client object")
	}

	header := make(http.Header)
	header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	header.Add("Content-Type", "application/json")
	header.Add("Accept", "application/json")

	return &Client{
		BaseURL:        u,
		BaseHttpHeader: header,
		httpClient:     http.Client{Jar: jar},
	}, nil
}

func (c *Client) Get(uri string, header http.Header) (*Response, error) {
	req, err := c.newRequest("GET", uri, nil, header)
	if err != nil {
		return nil, errors.WithMessage(err, "get request "+uri+" failed")
	}
	return c.do(req)
}

func (c *Client) Post(uri string, body interface{}, header http.Header) (*Response, error) {
	req, err := c.newRequest("POST", uri, body, header)
	if err != nil {
		return nil, errors.WithMessage(err, "post request "+uri+" failed")
	}
	return c.do(req)
}

func (c *Client) Put(uri string, body interface{}, header http.Header) (*Response, error) {
	req, err := c.newRequest("PUT", uri, body, header)
	if err != nil {
		return nil, errors.WithMessage(err, "put request "+uri+" failed")
	}
	return c.do(req)
}

func (c *Client) newRequest(method string, uri string, body interface{}, header http.Header) (*Request, error) {

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create request object")
		}
	}
	log.Printf("JSON: %v", buf)
	u := c.BaseURL.String() + uri
	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request object")
	}

	for k, v := range c.BaseHttpHeader {
		for _, h := range v {
			//fmt.Printf("Adding header (%s): (%s - %s)\n\n", u, k, h)
			req.Header.Add(k, h)
		}
	}

	for k, v := range header {
		for _, h := range v {
			//fmt.Printf("Adding header (%s): (%s - %s)\n\n", u, k, h)
			req.Header.Add(k, h)
		}
	}

	return &Request{
		Request: *req,
	}, nil
}

func (c *Client) newResponse(resp *http.Response) (*Response, error) {

	return &Response{
		Response: *resp,
	}, nil
}

func (c *Client) do(req *Request) (*Response, error) {

	//fmt.Printf("\n\nCOOKIES (%s): %v\n\n", req.URL, c.httpClient.Jar.Cookies(req.URL))
	//fmt.Printf("\n\nHEADERS (%s): %v\n\n", req.URL, req.Header)

	resp, err := c.httpClient.Do(&req.Request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute http request")
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, errors.New("http request failed with status: " + resp.Status)
	}

	return c.newResponse(resp)
}
