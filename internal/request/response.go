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
