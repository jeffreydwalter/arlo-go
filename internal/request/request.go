package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime"
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

type Request struct {
	http.Request
}
type Response struct {
	http.Response
	Data interface{}
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

func GetContentType(ct string) (string, error) {
	mediaType, _, err := mime.ParseMediaType(ct)

	if err != nil {
		return "", errors.Wrap(err, "failed to get content type")
	}
	return mediaType, nil
}

/*
func (resp *Response) Parse(schema interface{}) (interface{}, error){
	mediatype, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	log.Printf("CONTENT TYPE %s\n", mediatype)

	switch mediatype {
	case "application/json":
		log.Println("DECODING JSON: %s", json.Valid(resp.Data))
		if err := json.Unmarshal(resp.Data, schema); err != nil {
			log.Println("GOT AN ERROR")
			return nil, err
		}
	}

	return schema, nil
}
*/

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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// log.Printf("DATA: %v", string(data))

	var d interface{}
	mediaType, err := GetContentType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create response object")
	}

	switch mediaType {
	case "application/json":
		err = json.Unmarshal([]byte(data), &d)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create response object")
		}
	}

	return &Response{
		Response: *resp,
		Data:     d,
	}, nil
}

func (c *Client) do(req *Request) (*Response, error) {

	//fmt.Printf("\n\nCOOKIES (%s): %v\n\n", req.URL, c.httpClient.Jar.Cookies(req.URL))
	//fmt.Printf("\n\nHEADERS (%s): %v\n\n", req.URL, req.Header)

	resp, err := c.httpClient.Do(&req.Request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute http request")
	}
	defer resp.Body.Close()

	return c.newResponse(resp)
}
