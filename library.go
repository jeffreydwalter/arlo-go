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
	"time"

	"github.com/pkg/errors"
)

// LibraryMetaData is the library meta data.
type LibraryMetaData struct {
	DateTo   string                         `json:"dateTo"`
	DateFrom string                         `json:"dateFrom"`
	Meta     map[string]map[string]Favorite `json:"meta"`
}

// presignedContentUrl is a link to the actual video in Amazon AWS.
// presignedThumbnailUrl is a link to the thumbnail .jpg of the actual video in Amazon AWS.
type Recording struct {
	MediaDurationSecond   int    `json:"mediaDurationSecond"`
	ContentType           string `json:"contentType"`
	Name                  string `json:"name"`
	PresignedContentUrl   string `json:"presignedContentUrl"`
	LastModified          int64  `json:"lastModified"`
	LocalCreatedDate      int64  `json:"localCreatedDate"`
	PresignedThumbnailUrl string `json:"presignedThumbnailUrl"`
	Reason                string `json:"reason"`
	DeviceId              string `json:"deviceId"`
	CreatedBy             string `json:"createdBy"`
	CreatedDate           string `json:"createdDate"`
	TimeZone              string `json:"timeZone"`
	OwnerId               string `json:"ownerId"`
	UtcCreatedDate        int64  `json:"utcCreatedDate"`
	CurrentState          string `json:"currentState"`
	MediaDuration         string `json:"mediaDuration"`
	UniqueId              string `json:"uniqueId"`
}

type Library []Recording

func (a *Arlo) GetLibraryMetaData(fromDate, toDate time.Time) (libraryMetaData *LibraryMetaData, err error) {
	msg := "failed to get library metadata"

	body := map[string]string{"dateFrom": fromDate.Format("20060102"), "dateTo": toDate.Format("20060102")}
	resp, err := a.post(LibraryMetadataUri, "", body, nil)
	if err != nil {
		return nil, errors.WithMessage(err, msg)
	}
	defer resp.Body.Close()

	response := new(LibraryMetaDataResponse)
	if err := resp.Decode(&response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New(msg)
	}

	return &response.Data, nil
}

func (a *Arlo) GetLibrary(fromDate, toDate time.Time) (library *Library, err error) {
	msg := "failed to get library"

	body := map[string]string{"dateFrom": fromDate.Format("20060102"), "dateTo": toDate.Format("20060102")}
	resp, err := a.post(LibraryUri, "", body, nil)
	if err != nil {
		return nil, errors.WithMessage(err, msg)
	}
	defer resp.Body.Close()

	response := new(LibraryResponse)
	if err := resp.Decode(&response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New(msg)
	}

	return &response.Data, nil
}

/*
 Delete a single video recording from arlo.

 All of the date info and device id you need to pass into this method are given in the results of the GetLibrary() call.

 NOTE: {"data": [{"createdDate": r.CreatedDate, "utcCreatedDate": r.UtcCreatedDate, "deviceId": r.DeviceId}]} is all that's really required.
*/
func (a *Arlo) DeleteRecording(r *Recording) error {

	body := map[string]Library{"data": {*r}}
	resp, err := a.post(LibraryRecycleUri, "", body, nil)
	return checkRequest(resp, err, "failed to delete recording")
}

/*
 Delete a batch of video recordings from arlo.

 The GetLibrary() call response json can be passed directly to this method if you'd like to delete the same list of videos you queried for.

 NOTE: {"data": [{"createdDate": r.CreatedDate, "utcCreatedDate": r.UtcCreatedDate, "deviceId": r.DeviceId}]} is all that's really required.
*/
func (a *Arlo) BatchDeleteRecordings(l *Library) error {

	body := map[string]Library{"data": *l}
	resp, err := a.post(LibraryRecycleUri, "", body, nil)
	return checkRequest(resp, err, "failed to delete recordings")
}
