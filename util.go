package arlo

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jeffreydwalter/arlo-golang/internal/request"
	"github.com/jeffreydwalter/arlo-golang/internal/util"

	"github.com/pkg/errors"
)

func checkHttpRequest(resp *request.Response, err error, msg string) error {
	if resp.StatusCode != 200 {
		return errors.WithMessage(errors.New(fmt.Sprintf("http request failed: %s (%d)", resp.Status, resp.StatusCode)), msg)
	}

	if err != nil {
		return errors.WithMessage(err, msg)
	}

	return nil
}

func checkRequest(resp *request.Response, err error, msg string) error {
	defer resp.Body.Close()

	if err := checkHttpRequest(resp, err, msg); err != nil {
		return err
	}

	var status Status
	if err := resp.Decode(&status); err != nil {
		return err
	}

	if status.Success == false {
		return errors.WithMessage(errors.New(status.Reason), msg)
	}

	return nil
}

func genTransId() string {

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	e := random.Float64() * math.Pow(2, 32)

	ms := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

	return fmt.Sprintf("%s!%s!%s", TransIdPrefix, strings.ToLower(util.FloatToHex(e)), strconv.Itoa(int(ms)))
}

func (a *Arlo) get(uri, xCloudId string, header http.Header) (*request.Response, error) {
	if len(xCloudId) > 0 {
		a.rwmutex.Lock()
		a.client.BaseHttpHeader.Set("xcloudId", xCloudId)
		a.rwmutex.Unlock()
	}

	return a.client.Get(uri, header)
}

func (a *Arlo) put(uri, xCloudId string, body interface{}, header http.Header) (*request.Response, error) {
	if len(xCloudId) > 0 {
		a.rwmutex.Lock()
		a.client.BaseHttpHeader.Set("xcloudId", xCloudId)
		a.rwmutex.Unlock()
	}

	return a.client.Put(uri, body, header)
}

func (a *Arlo) post(uri, xCloudId string, body interface{}, header http.Header) (*request.Response, error) {
	if len(xCloudId) > 0 {
		a.rwmutex.Lock()
		a.client.BaseHttpHeader.Set("xcloudId", xCloudId)
		a.rwmutex.Unlock()
	}

	return a.client.Post(uri, body, header)
}
