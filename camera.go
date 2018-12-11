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
	"strings"
	"time"

	"github.com/pkg/errors"
)

// A Camera is a Device of type "camera".
// This type is here just for semantics. Some methods explicitly require a device of a certain type.
type Camera Device

// Cameras is a slice of Camera objects.
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

// PushToTalk starts a push-to-talk session.
// FIXME: This feature requires more API calls to make it actually work, and I haven't figure out how to fully implement it.
// It appears that the audio stream is Real-Time Transport Protocol (RTP), which requires a player (ffmpeg?) to consume the stream.
func (c *Camera) PushToTalk() error {
	/*
		processResponse: function(e) {
		            if (g.pc)
		                if (e.properties && "answerSdp" == e.properties.type) {
		                    var t = e.properties.data
		                      , i = {
		                        type: "answer",
		                        sdp: t
		                    };
		                    r.debug(i),
		                    g.pc.setRemoteDescription(new g.SessionDescription(i), u, d)
		                } else if (e.properties && "answerCandidate" == e.properties.type)
		                    if (g.candidateCache)
		                        g.candidateCache.push(e.properties.data);
		                    else {
		                        var n = e.properties.data
		                          , a = window.mozRTCIceCandidate || window.RTCIceCandidate
		                          , o = new a({
		                            candidate: n,
		                            sdpMLineIndex: 0
		                        });
		                        r.debug(o),
		                        g.pc.addIceCandidate(o)
		                    }
		        },
		        startConnection: function(t) {
		            g.loading = !0,
		            g.error = !1,
		            g.candidateCache = [];
		            var i = t.deviceId
		              , o = t.parentId
		              , u = t.uniqueId;
		            g.device = t;
		            var p = {
		                method: "GET",
		                url: l.getPttUrl(u),
		                data: "",
		                headers: {
		                    Authorization: s.ssoToken,
		                    "Content-Type": "application/json; charset=utf-8",
		                    "Data-Type": "json"
		                }
		            };
		            r.debug("getting ptt data: " + JSON.stringify(p));
		            n(p).then(function(u) {
		                if (!u.data.success)
		                    return e.$broadcast("show_error", u.data),
		                    void (g.error = u.data.data.message || !0);
		                var m = u.data.data.data;
		                g.uSessionId = u.data.data.uSessionId,
		                _.each(m, function(e) {
		                    e.url && (e.urls = e.url,
		                    delete e.url)
		                });
		                var f = new g.PeerConnection({
		                    iceServers: m,
		                    iceCandidatePoolSize: 0
		                });
		                f.onicecandidate = function(e) {
		                    if (null != e.candidate) {
		                        r.debug(e.candidate);
		                        var a = {
		                            action: "pushToTalk",
		                            from: t.userId,
		                            publishResponse: !1,
		                            resource: "cameras/" + i,
		                            responseUrl: "",
		                            to: o,
		                            transId: "web!98b0c88b!1429756137177",
		                            properties: {
		                                uSessionId: g.uSessionId,
		                                type: "offerCandidate",
		                                data: e.candidate.candidate
		                            }
		                        };
		                        p = {
		                            method: "POST",
		                            url: l.getPttNotifyUrl(o),
		                            data: a,
		                            headers: {
		                                xcloudId: t.xCloudId,
		                                Authorization: s.ssoToken
		                            }
		                        },
		                        n(p)
		                    } else
		                        r.debug("Failed to get any more candidate")
		                }
		                ,
		                f.oniceconnectionstatechange = function(e) {
		                    r.debug("ICE Connection State Change:" + f.iceConnectionState),
		                    "connected" == f.iceConnectionState || "completed" == f.iceConnectionState ? g.loading = !1 : "disconnected" != f.iceConnectionState && "failed" != f.iceConnectionState || (g.stopConnection(),
		                    g.error = a("i18n")("camera_label_ptt_failed_to_connect"))
		                }
		                ,
		                g.pc = f,
		                (navigator.getUserMedia || navigator.webkitGetUserMedia || navigator.mozGetUserMedia || navigator.msGetUserMedia).call(navigator, {
		                    audio: !0,
		                    video: !1
		                }, function(e) {
		                    g.stream = e,
		                    g.stream.getAudioTracks()[0].enabled = !1,
		                    f.addStream(e),
		                    f.createOffer(function(e) {
		                        f.setLocalDescription(e, c, d),
		                        r.debug(e.sdp);
		                        var a = {
		                            action: "pushToTalk",
		                            from: t.userId,
		                            publishResponse: !0,
		                            resource: "cameras/" + i,
		                            responseUrl: "",
		                            to: o,
		                            transId: "web!98b0c88b!1429756137177",
		                            properties: {
		                                uSessionId: g.uSessionId,
		                                type: "offerSdp",
		                                data: e.sdp
		                            }
		                        };
		                        p = {
		                            method: "POST",
		                            url: l.getPttNotifyUrl(o),
		                            data: a,
		                            headers: {
		                                xcloudId: t.xCloudId,
		                                Authorization: s.ssoToken
		                            }
		                        },
		                        n(p)
		                    }, d)
		                }, d)
		            })
		        },
		        stopConnection: function() {
		            if (g.pc) {
		                var e = {
		                    action: "pushToTalk",
		                    from: g.device.userId,
		                    publishResponse: !1,
		                    resource: "cameras/" + g.device.deviceId,
		                    responseUrl: "",
		                    to: g.device.deviceId,
		                    transId: "web!98b0c88b!1429756137177",
		                    properties: {
		                        uSessionId: g.uSessionId,
		                        type: "endSession"
		                    }
		                }
		                  , t = {
		                    method: "POST",
		                    url: l.getPttNotifyUrl(g.device.deviceId),
		                    data: e,
		                    headers: {
		                        xcloudId: g.device.xCloudId,
		                        Authorization: s.ssoToken
		                    }
		                };
		                n(t);
		                try {
		                    g.stream.getAudioTracks()[0].stop(),
		                    g.stream = null
		                } catch (e) {}
		                g.pc.close(),
		                g.pc = null,
		                g.loading = !0
		            }
		        }
		    };
	*/
	resp, err := c.arlo.get(fmt.Sprintf(PttUri, c.UniqueId), c.XCloudId, nil)
	return checkRequest(resp, err, "failed to enable push to talk")
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

// StartStream returns a json object containing the rtmps url to the requested video stream.
// You will need something like ffmpeg to read the rtmps stream.

// If you call StartStream(), you have to start reading data from the stream, or streaming will be cancelled
// and taking a snapshot may fail (since it requires the stream to be active).
func (c *Camera) StartStream() (url string, err error) {
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

	resp, err := c.arlo.post(StartStreamUri, c.XCloudId, payload, nil)
	if err != nil {
		return "", errors.WithMessage(err, msg)
	}
	defer resp.Body.Close()

	response := new(StreamResponse)
	if err := resp.Decode(response); err != nil {
		return "", err
	}

	if !response.Success {
		return "", errors.WithMessage(errors.New("status was false"), msg)
	}

	response.Data.URL = strings.Replace(response.Data.URL, "rtsp://", "rtsps://", 1)

	return response.Data.URL, nil
}

// TakeSnapshot causes the camera to snapshot while recording.
// NOTE: You MUST call StartStream() before calling this function.
// If you call StartStream(), you have to start reading data from the stream, or streaming will be cancelled
// and taking a snapshot may fail (since it requires the stream to be active).

// NOTE: You should not use this function is you just want a snapshot and aren't intending to stream.
// Use TriggerFullFrameSnapshot() instead.
//
// NOTE: Use DownloadSnapshot() to download the actual image file.
// TODO: Need to refactor the even stream code to allow handling of events whose transIds don't correlate. :/
func (c *Camera) TakeSnapshot() (response *EventStreamResponse, err error) {

	return nil, errors.New("TakeSnapshot not implemented")
	/*
		msg := "failed to take snapshot"

		body := map[string]string{"deviceId": c.DeviceId, "parentId": c.ParentId, "xcloudId": c.XCloudId, "olsonTimeZone": c.Properties.OlsonTimeZone}
		resp, err := c.arlo.post(TakeSnapshotUri, c.XCloudId, body, nil)
		if err := checkRequest(resp, err, msg); err != nil {
			return nil, errors.WithMessage(err, msg)
		}
	*/

	// TODO: Need to write the code to handle the event stream message.
	/*
			def callback(self, event):
				if event.get("deviceId") == camera.get("deviceId") and event.get("resource") == "mediaUploadNotification":
					presigned_content_url = event.get("presignedContentUrl")
					if presigned_content_url is not None:
		r				return presigned_content_url
	*/
}

// TriggerFullFrameSnapshot causes the camera to record a full-frame snapshot.
// The presignedFullFrameSnapshotUrl url is returned.
// Use DownloadSnapshot() to download the actual image file.
// TODO: Need to refactor the even stream code to allow handling of events whose transIds don't correlate. :/
func (c *Camera) TriggerFullFrameSnapshot() (response *EventStreamResponse, err error) {

	return nil, errors.New("TriggerFullFrameSnapshot not implemented")
	/*
		payload := EventStreamPayload{
			Action:          "set",
			Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
			PublishResponse: true,
			Properties: map[string]string{
				"activityState": "fullFrameSnapshot",
			},
			TransId: genTransId(),
			From:    fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
			To:      c.ParentId,
		}

		msg := "failed to trigger full-frame snapshot"

		b := c.arlo.Basestations.Find(c.ParentId)
		if b == nil {
			err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
			return nil, errors.WithMessage(err, msg)
		}
		return b.makeEventStreamRequest(payload, msg)
	*/
	/*
		def callback(self, event):
			if event.get("from") == basestation.get("deviceId") and event.get("resource") == "cameras/"+camera.get("deviceId") and event.get("action") == "fullFrameSnapshotAvailable":
				return event.get("properties", {}).get("presignedFullFrameSnapshotUrl")
			return None
	*/
}

// StartRecording causes the camera to start recording and returns a url that you must start reading from using ffmpeg
// or something similar.
func (c *Camera) StartRecording() (url string, err error) {
	msg := "failed to start recording"

	url, err = c.StartStream()
	if err != nil {
		return "", errors.WithMessage(err, msg)
	}

	body := map[string]string{"deviceId": c.DeviceId, "parentId": c.ParentId, "xcloudId": c.XCloudId, "olsonTimeZone": c.Properties.OlsonTimeZone}
	resp, err := c.arlo.post(StartRecordUri, c.XCloudId, body, nil)
	if err := checkRequest(resp, err, msg); err != nil {
		return "", errors.WithMessage(err, msg)
	}

	return url, nil
}

// StopRecording causes the camera to stop recording.
func (c *Camera) StopRecording() error {
	msg := "failed to stop recording"

	body := map[string]string{"deviceId": c.DeviceId, "parentId": c.ParentId, "xcloudId": c.XCloudId, "olsonTimeZone": c.Properties.OlsonTimeZone}
	resp, err := c.arlo.post(StopRecordUri, c.XCloudId, body, nil)
	if err := checkRequest(resp, err, msg); err != nil {
		return errors.WithMessage(err, msg)
	}

	return nil
}

// This function downloads a Cvr Playlist file for the period fromDate to toDate.
func (c *Camera) GetCvrPlaylist(fromDate, toDate time.Time) (playlist *CvrPlaylist, err error) {
	msg := "failed to get cvr playlist"

	resp, err := c.arlo.get(fmt.Sprintf(PlaylistUri, c.UniqueId, fromDate.Format("20060102"), toDate.Format("20060102")), c.XCloudId, nil)

	if err != nil {
		return nil, errors.WithMessage(err, msg)
	}
	defer resp.Body.Close()

	response := new(CvrPlaylistResponse)
	if err := resp.Decode(&response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New(msg)
	}

	return &response.Data, nil
}
