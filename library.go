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
	resp, err := a.post(MetadataUri, "", body, nil)
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
	resp, err := a.post(RecordingsUri, "", body, nil)
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
	resp, err := a.post(RecycleUri, "", body, nil)
	return checkRequest(resp, err, "failed to delete recording")
}

/*
 Delete a batch of video recordings from arlo.

 The GetLibrary() call response json can be passed directly to this method if you'd like to delete the same list of videos you queried for.

 NOTE: {"data": [{"createdDate": r.CreatedDate, "utcCreatedDate": r.UtcCreatedDate, "deviceId": r.DeviceId}]} is all that's really required.
*/
func (a *Arlo) BatchDeleteRecordings(l *Library) error {
	body := map[string]Library{"data": *l}
	resp, err := a.post(RecycleUri, "", body, nil)
	return checkRequest(resp, err, "failed to delete recordings")
}

// SendAnalyticFeedback is only really used by the GUI. It is a response to a prompt asking you whether an object which
// was tagged by it's AI in your recording was tagged correctly.
func (a *Arlo) SendAnalyticFeedback(r *Recording) error {
	category := "Person" // Other
	body := map[string]map[string]interface{}{"data": {"utcCreatedDate": r.UtcCreatedDate, "category": category, "createdDate": r.CreatedDate}}
	resp, err := a.put(AnalyticFeedbackUri, "", body, nil)
	return checkRequest(resp, err, "failed to send analytic feedback about recording")
}

// GetActiveAutomationDefinitions gets the mode metadata (this API replaces the older GetModes(), which still works).
func (a *Arlo) GetActiveAutomationDefinitions() error {
	resp, err := a.get(ActiveAutomationUri, "", nil)
	return checkRequest(resp, err, "failed to get active automation definitions")
}

/*
func (a *Arlo) SetActiveAutomationMode() error {

	body := struct{}{} //map[string]map[string]interface{}{"data": {"utcCreatedDate": r.UtcCreatedDate, "category": category, "createdDate": r.CreatedDate}}
	resp, err := a.put(AnalyticFeedbackUri, "", body, nil)
	return checkRequest(resp, err, "failed to send analytic feedback about recording")
}
*/
/*
[
    {
        "activeModes": [
            "mode1"
        ],
        "activeSchedules": [],
        "gatewayId": "48935B7SA9847",
        "schemaVersion": 1,
        "timestamp": 1536781758034,
        "type": "activeAutomations",
        "uniqueId": "336-4764296_48935B7SA9847"
    }
]
*/
/*
   setActiveAutomationMode: function(r, a) {
       var s = {
           activeAutomations: [{
               deviceId: a.gatewayId,
               timestamp: _.now(),
               activeModes: [r],
               activeSchedules: []
           }]
       }
         , l = {
           method: "POST",
           data: s,
           url: d.getActiveAutomationUrl(a.gatewayId),
           headers: {
               Authorization: o.ssoToken,
               schemaVersion: 1
           }
       };
       return n.debug("calling set active automation mode with config:" + JSON.stringify(l)),
       i(l).then(function(i) {
           if (n.debug("got set active automation mode result:" + JSON.stringify(i)),
           i && i.data && !i.data.success)
               return e.$broadcast(c.appEvents.SHOW_ERROR, i.data),
               t.reject(i.data)
       })
   },
   setActiveAutomationSchedule: function(r) {
       var r = {
           activeAutomations: [{
               deviceId: r.deviceId,
               timestamp: _.now(),
               activeModes: [],
               activeSchedules: [r.scheduleId]
           }]
       }
         , a = {
           method: "POST",
           data: r,
           url: d.getActiveAutomationUrl(r.deviceId),
           headers: {
               Authorization: o.ssoToken,
               schemaVersion: 1
           }
       }
         , s = this;
       return n.debug("calling set active automation schedule with config:" + JSON.stringify(a)),
       i(a).then(function(i) {
           return n.debug("got set active automation schedule result:" + JSON.stringify(i)),
           i && i.data && !i.data.success ? (e.$broadcast(c.appEvents.SHOW_ERROR, i.data),
           t.reject(i.data)) : i && i.data && i.data.success ? (_.filter(s.activeAutomationDefinitions, function(e) {
               e.gatewayId == i.config.data.activeAutomations[0].deviceId && (e.activeModes = i.config.data.activeAutomations[0].activeModes,
               e.activeSchedules = i.config.data.activeAutomations[0].activeSchedules)
           }),
           i.config.data.activeAutomations[0]) : void 0
       })
   },
*/
