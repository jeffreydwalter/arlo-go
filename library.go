package arlo

import (
	"time"
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
}

type Library []Recording

func (a *Arlo) GetLibraryMetaData(fromDate, toDate time.Time) (response *LibraryMetaDataResponse, err error) {

	body := map[string]string{"dateFrom": fromDate.Format("20060102"), "dateTo": toDate.Format("20060102")}
	resp, err := a.post(LibraryMetadataUri, "", body, nil)
	if err := checkHttpRequest(resp, err, "failed to get library metadata"); err != nil {
		return nil, err
	}

	if err := resp.Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (a *Arlo) GetLibrary(fromDate, toDate time.Time) (response *LibraryResponse, err error) {

	body := map[string]string{"dateFrom": fromDate.Format("20060102"), "dateTo": toDate.Format("20060102")}
	resp, err := a.post(LibraryUri, "", body, nil)
	if err := checkHttpRequest(resp, err, "failed to get library"); err != nil {
		return nil, err
	}

	if err := resp.Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}

/*
 Delete a single video recording from arlo.

 All of the date info and device id you need to pass into this method are given in the results of the GetLibrary() call.

 NOTE: {"data": [{"createdDate": r.CreatedDate, "utcCreatedDate": r.UtcCreatedDate, "deviceId": r.DeviceId}]} is all that's really required.
*/
func (a *Arlo) DeleteRecording(r Recording) error {

	body := map[string]Library{"data": {r}}
	resp, err := a.post(LibraryRecycleUri, "", body, nil)
	return checkRequest(resp, err, "failed to delete recording")
}

/*
 Delete a batch of video recordings from arlo.

 The GetLibrary() call response json can be passed directly to this method if you'd like to delete the same list of videos you queried for.

 NOTE: {"data": [{"createdDate": r.CreatedDate, "utcCreatedDate": r.UtcCreatedDate, "deviceId": r.DeviceId}]} is all that's really required.
*/
func (a *Arlo) BatchDeleteRecordings(l Library) error {

	body := map[string]Library{"data": l}
	resp, err := a.post(LibraryRecycleUri, "", body, nil)
	return checkRequest(resp, err, "failed to delete recordings")
}
