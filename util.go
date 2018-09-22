/*
 * Copyright (c) 2018 Jeffrey Walter <jeffreydwalter@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package arlo

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jeffreydwalter/arlo-golang/internal/request"
	"github.com/jeffreydwalter/arlo-golang/internal/util"

	"github.com/pkg/errors"
)

func checkRequest(resp *request.Response, err error, msg string) error {
	if err != nil {
		return errors.WithMessage(err, msg)
	}
	defer resp.Body.Close()

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
	a.client.AddHeader("xcloudId", xCloudId)
	return a.client.Get(uri, header)
}

func (a *Arlo) put(uri, xCloudId string, body interface{}, header http.Header) (*request.Response, error) {
	a.client.AddHeader("xcloudId", xCloudId)
	return a.client.Put(uri, body, header)
}

func (a *Arlo) post(uri, xCloudId string, body interface{}, header http.Header) (*request.Response, error) {
	a.client.AddHeader("xcloudId", xCloudId)
	return a.client.Post(uri, body, header)
}

func (a *Arlo) DownloadFile(url, to string) error {
	msg := fmt.Sprintf("failed to download file (%s) => (%s)", url, to)
	resp, err := a.get(url, "", nil)
	if err != nil {
		return errors.WithMessage(err, msg)
	}
	defer resp.Body.Close()

	f, err := os.Create(to)
	if err != nil {
		return errors.WithMessage(err, msg)
	}

	_, err = io.Copy(f, resp.Body)
	defer f.Close()
	if err != nil {
		return errors.WithMessage(err, msg)
	}

	return nil
}
