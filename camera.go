package arlo

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// A Camera is a Device of type "camera".
// This type is here just for semantics. Some methods explicitly require a device of a certain type.
type Camera Device

// Cameras is an array of Camera objects.
type Cameras []Camera

// StartStream returns a json object containing the rtmps url to the requested video stream.
// You will need the to install a library to handle streaming of this protocol: https://pypi.python.org/pypi/python-librtmp
//
// The request to /users/devices/startStream returns:
// NOTE: { "url":"rtsp://vzwow09-z2-prod.vz.netgear.com:80/vzmodulelive?egressToken=b1b4b675_ac03_4182_9844_043e02a44f71&userAgent=web&cameraId=48B4597VD8FF5_1473010750131" }
func (c *Camera) StartStream() (*StreamResponse, error) {
	body := Payload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: map[string]string{
			"activityState": "startUserStream",
			"cameraId":      c.DeviceId,
		},
		TransId: genTransId(),
		From:    fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:      c.ParentId,
	}

	resp, err := c.arlo.post(DeviceStartStreamUri, c.XCloudId, body, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to start stream")
	}
	defer resp.Body.Close()

	var streamResponse StreamResponse
	if err := resp.Decode(&streamResponse); err != nil {
		return nil, err
	}

	if !streamResponse.Success {
		return nil, errors.WithMessage(errors.New("status was false"), "failed to start stream")
	}

	streamResponse.Data.Url = strings.Replace(streamResponse.Data.Url, "rtsp://", "rtsps://", 1)

	return &streamResponse, nil
}

// TakeSnapshot causes the camera to record a snapshot.
func (c *Camera) TakeSnapshot() (*StreamResponse, error) {
	streamResponse, err := c.StartStream()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to take snapshot")
	}

	body := map[string]string{"deviceId": c.DeviceId, "parentId": c.ParentId, "xcloudId": c.XCloudId, "olsonTimeZone": c.Properties.OlsonTimeZone}
	resp, err := c.arlo.post(DeviceTakeSnapshotUri, c.XCloudId, body, nil)
	if err := checkRequest(*resp, err, "failed to update device name"); err != nil {
		return nil, errors.WithMessage(err, "failed to take snapshot")
	}

	return streamResponse, nil
}

// StartRecording causes the camera to start recording and returns a url that you must start reading from using ffmpeg
// or something similar.
func (c *Camera) StartRecording() (*StreamResponse, error) {
	streamResponse, err := c.StartStream()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to start recording")
	}

	body := map[string]string{"deviceId": c.DeviceId, "parentId": c.ParentId, "xcloudId": c.XCloudId, "olsonTimeZone": c.Properties.OlsonTimeZone}
	resp, err := c.arlo.post(DeviceStartRecordUri, c.XCloudId, body, nil)
	if err := checkRequest(*resp, err, "failed to update device name"); err != nil {
		return nil, errors.WithMessage(err, "failed to start recording")
	}

	return streamResponse, nil
}
