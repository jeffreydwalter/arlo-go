package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/pkg/errors"

	"golang.org/x/net/publicsuffix"
)

type Client struct {
	BaseURL        *url.URL
	BaseHttpHeader *http.Header
	HttpClient     *http.Client
}

func NewClient(baseurl string) (*Client, error) {
	var err error
	var jar *cookiejar.Jar

	options := cookiejar.Options{PublicSuffixList: publicsuffix.List}

	if jar, err = cookiejar.New(&options); err != nil {
		return nil, errors.Wrap(err, "failed to create client object")
	}

	var u *url.URL
	if u, err = url.Parse(baseurl); err != nil {
		return nil, errors.Wrap(err, "failed to create client object")
	}

	header := make(http.Header)
	header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_1_2 like Mac OS X) AppleWebKit/604.3.5 (KHTML, like Gecko) Mobile/15B202 NETGEAR/v1 (iOS Vuezone)")
	header.Set("Content-Type", "application/json")
	header.Set("Accept", "application/json")

	return &Client{
		BaseURL:        u,
		BaseHttpHeader: &header,
		HttpClient:     &http.Client{Jar: jar},
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
	// log.Printf("\n\nBODY (%s): %s\n\n", uri, buf)

	u := c.BaseURL.String() + uri
	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request object")
	}

	for k, v := range *c.BaseHttpHeader {
		for _, h := range v {
			//log.Printf("Adding header (%s): (%s - %s)\n\n", u, k, h)
			req.Header.Set(k, h)
		}
	}

	for k, v := range header {
		for _, h := range v {
			//log.Printf("Adding header (%s): (%s - %s)\n\n", u, k, h)
			req.Header.Set(k, h)
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

	//log.Printf("\n\nCOOKIES (%s): %v\n\n", req.URL, c.HttpClient.Jar.Cookies(req.URL))
	//log.Printf("\n\nHEADERS (%s): %v\n\n", req.URL, req.Header)

	resp, err := c.HttpClient.Do(&req.Request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute http request")
	}

	if resp.StatusCode >= http.StatusBadRequest {
		defer resp.Body.Close()
		return nil, errors.New("http request failed with status: " + resp.Status)
	}

	return c.newResponse(resp)
}
