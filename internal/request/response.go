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

package request

import (
	"encoding/json"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"reflect"

	"github.com/pkg/errors"
)

type Response struct {
	http.Response
}

func (resp *Response) GetContentType() (string, error) {

	mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return "", errors.Wrap(err, "failed to get content type")
	}
	return mediaType, nil
}

func (resp *Response) Decode(s interface{}) error {

	defer resp.Body.Close()

	mediaType, err := resp.GetContentType()
	if err != nil {
		return errors.WithMessage(err, "failed to decode response body")
	}

	switch mediaType {
	case "application/json":
		err := json.NewDecoder(resp.Body).Decode(&s)
		if err != nil {
			return errors.Wrap(err, "failed to create "+reflect.TypeOf(s).String()+" object")
		}
	default:
		return errors.New("unsupported content type: " + mediaType)
	}
	return nil
}

func (resp *Response) Download(to string) (error, int64) {

	defer resp.Body.Close()

	// Create output file
	newFile, err := os.Create(to)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	// Write bytes from HTTP response to file.
	// response.Body satisfies the reader interface.
	// newFile satisfies the writer interface.
	// That allows us to use io.Copy which accepts
	// any type that implements reader and writer interface
	bytesWritten, err := io.Copy(newFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return nil, bytesWritten
}
