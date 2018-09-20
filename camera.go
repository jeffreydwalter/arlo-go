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

// Find returns a camera with the device id passed in.
func (cs *Cameras) Find(deviceId string) *Camera {
	for _, c := range *cs {
		if c.DeviceId == deviceId {
			return &c
		}
	}

	return nil
}

// On turns a camera on; meaning it will detect and record events.
func (c *Camera) On() (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: CameraProperties{
			PrivacyActive: false,
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to turn camera on"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

// On turns a camera off; meaning it won't detect and record events.
func (c *Camera) Off() (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: CameraProperties{
			PrivacyActive: true,
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to turn camera off"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

// SetBrightness sets the camera brightness.
// NOTE: Brightness is between -2 and 2 in increments of 1 (-2, -1, 0, 1, 2).
// Setting it to an invalid value has no effect.
func (c *Camera) SetBrightness(brightness int) (response *EventStreamResponse, err error) {
	// Sanity check; if the values are above or below the allowed limits, set them to their limit.
	if brightness < -2 {
		brightness = -2
	} else if brightness > 2 {
		brightness = 2
	}

	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: CameraProperties{
			Brightness: brightness,
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to set camera brightness"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

// StartStream returns a json object containing the rtmps url to the requested video stream.
// You will need the to install a library to handle streaming of this protocol: https://pypi.python.org/pypi/python-librtmp
//
// The request to /users/devices/startStream returns:
// NOTE: { "url":"rtsp://vzwow09-z2-prod.vz.netgear.com:80/vzmodulelive?egressToken=b1b4b675_ac03_4182_9844_043e02a44f71&userAgent=web&cameraId=48B4597VD8FF5_1473010750131" }
func (c *Camera) StartStream() (response *StreamResponse, err error) {
	payload := EventStreamPayload{
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

	msg := "failed to start stream"

	resp, err := c.arlo.post(DeviceStartStreamUri, c.XCloudId, payload, nil)
	if err := checkHttpRequest(resp, err, msg); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := resp.Decode(response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.WithMessage(errors.New("status was false"), msg)
	}

	response.Data.Url = strings.Replace(response.Data.Url, "rtsp://", "rtsps://", 1)

	return response, nil
}

// TakeSnapshot causes the camera to record a snapshot.
func (c *Camera) TakeSnapshot() (response *StreamResponse, err error) {
	msg := "failed to take snapshot"

	response, err = c.StartStream()
	if err != nil {
		return nil, errors.WithMessage(err, msg)
	}

	body := map[string]string{"deviceId": c.DeviceId, "parentId": c.ParentId, "xcloudId": c.XCloudId, "olsonTimeZone": c.Properties.OlsonTimeZone}
	resp, err := c.arlo.post(DeviceTakeSnapshotUri, c.XCloudId, body, nil)
	if err := checkRequest(resp, err, "failed to update device name"); err != nil {
		return nil, errors.WithMessage(err, msg)
	}

	return response, nil
}

// StartRecording causes the camera to start recording and returns a url that you must start reading from using ffmpeg
// or something similar.
func (c *Camera) StartRecording() (response *StreamResponse, err error) {
	msg := "failed to start recording"

	response, err = c.StartStream()
	if err != nil {
		return nil, errors.WithMessage(err, msg)
	}

	body := map[string]string{"deviceId": c.DeviceId, "parentId": c.ParentId, "xcloudId": c.XCloudId, "olsonTimeZone": c.Properties.OlsonTimeZone}
	resp, err := c.arlo.post(DeviceStartRecordUri, c.XCloudId, body, nil)
	if err := checkRequest(resp, err, "failed to update device name"); err != nil {
		return nil, errors.WithMessage(err, msg)
	}

	return response, nil
}

func (c *Camera) EnableMotionAlerts(sensitivity int, zones []string) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: MotionDetectionProperties{
			BaseDetectionProperties: BaseDetectionProperties{
				Armed:       true,
				Sensitivity: sensitivity,
				Zones:       zones,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to enable motion alerts"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

func (c *Camera) DisableMotionAlerts(sensitivity int, zones []string) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: MotionDetectionProperties{
			BaseDetectionProperties: BaseDetectionProperties{
				Armed:       false,
				Sensitivity: sensitivity,
				Zones:       zones,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to enable motion alerts"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

func (c *Camera) EnableAudioAlerts(sensitivity int) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: AudioDetectionProperties{
			BaseDetectionProperties: BaseDetectionProperties{
				Armed:       true,
				Sensitivity: sensitivity,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to enable audio alerts"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

func (c *Camera) DisableAudioAlerts(sensitivity int) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: AudioDetectionProperties{
			BaseDetectionProperties: BaseDetectionProperties{
				Armed:       false,
				Sensitivity: sensitivity,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to disable audio alerts"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

// action: disabled OR recordSnapshot OR recordVideo
func (c *Camera) SetAlertNotificationMethods(action string, email, push bool) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: EventActionProperties{
			BaseEventActionProperties: BaseEventActionProperties{
				ActionType: action,
				StopType:   "timeout",
				Timeout:    15,
				EmailNotification: EmailNotification{
					Enabled:          email,
					EmailList:        []string{"__OWNER_EMAIL__"},
					PushNotification: push,
				},
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to set alert notification methods"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}
