package arlo

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

const eventStreamTimeout = 10 * time.Second
const pingTime = 30 * time.Second

// A Basestation is a Device that's not type "camera" (basestation, arloq, arloqs, etc.).
// This type is here just for semantics. Some methods explicitly require a device of a certain type.
type Basestation struct {
	Device
	eventStream *eventStream
}

// Basestations is an array of Basestation objects.
type Basestations []Basestation

// Find returns a basestation with the device id passed in.
func (bs *Basestations) Find(deviceId string) *Basestation {
	for _, b := range *bs {
		if b.DeviceId == deviceId {
			return &b
		}
	}

	return nil
}

// makeEventStreamRequest is a helper function sets up a response channel, sends a message to the event stream, and blocks waiting for the response.
func (b *Basestation) makeEventStreamRequest(payload EventStreamPayload, msg string) (response *EventStreamResponse, err error) {
	transId := genTransId()
	payload.TransId = transId

	if err := b.IsConnected(); err != nil {
		return nil, errors.WithMessage(err, msg)
	}

	subscriber := make(subscriber)

	// Add the response channel to the event stream queue so the response can be written to it.
	b.eventStream.subscribe(transId, subscriber)
	// Make sure we close and remove the response channel before returning.
	defer b.eventStream.unsubscribe(transId)

	// Send the payload to the event stream.
	if err := b.NotifyEventStream(payload, msg); err != nil {
		return nil, err
	}

	timer := time.NewTimer(eventStreamTimeout)
	defer timer.Stop()

	// Wait for the response to come back from the event stream on the response channel.
	select {
	// If we get a response, return it to the caller.
	case response := <-subscriber:
		return response, nil
	case err = <-b.eventStream.Error:
		return nil, errors.Wrap(err, msg)
	// If the event stream is closed, return an error about it.
	case <-b.eventStream.Disconnected:
		err = errors.New("event stream was closed before response was read")
		return nil, errors.WithMessage(err, msg)
	// If we timeout, return an error about it.
	case <-timer.C:
		err = fmt.Errorf("event stream response timed out after %.0f second", eventStreamTimeout.Seconds())
		return nil, errors.WithMessage(err, msg)
	}
}

func (b *Basestation) IsConnected() error {
	// If the event stream is closed, return an error about it.
	select {
	case <-b.eventStream.Disconnected:
		return errors.New("basestation not connected to event stream")
	default:
		return nil
	}
}

func (b *Basestation) Subscribe() error {
	b.eventStream = newEventStream(BaseUrl+fmt.Sprintf(SubscribeUri, b.arlo.Account.Token), b.arlo.client.HttpClient)

forLoop:
	for {
		// We blocking here because we can't really do anything with the event stream until we're connected.
		// Once we have confirmation that we're connected to the event stream, we will "subscribe" to events.
		select {
		case connected := <-b.eventStream.listen():
			if connected {
				break forLoop
			} else {
				return errors.New("failed to subscribe to the event stream")
			}
		case <-b.eventStream.Disconnected:
			err := errors.New("event stream was closed")
			return errors.WithMessage(err, "failed to subscribe to the event stream")
		}
	}

	if err := b.Ping(); err != nil {
		return errors.WithMessage(err, "failed to subscribe to the event stream")
	}

	// The Arlo event stream requires a "ping" every 30s.
	go func() {
		for {
			time.Sleep(pingTime)
			if err := b.Ping(); err != nil {
				b.Disconnect()
				break
			}
		}
	}()

	return nil
}

func (b *Basestation) Disconnect() error {
	// disconnect channel to stop event stream.
	if b.eventStream != nil {
		b.eventStream.disconnect()
	}
	return nil
}

// Ping makes a call to the subscriptions endpoint. The Arlo event stream requires this message to be sent every 30s.
func (b *Basestation) Ping() error {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("subscriptions/%s_%s", b.UserId, TransIdPrefix),
		PublishResponse: false,
		Properties:      map[string][1]string{"devices": {b.DeviceId}},
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	if _, err := b.makeEventStreamRequest(payload, "failed to ping the event stream"); err != nil {
		return err
	}
	return nil
}

func (b *Basestation) NotifyEventStream(payload EventStreamPayload, msg string) error {
	resp, err := b.arlo.post(fmt.Sprintf(NotifyUri, b.DeviceId), b.XCloudId, payload, nil)
	if err := checkRequest(resp, err, msg); err != nil {
		return errors.WithMessage(err, "failed to notify event stream")
	}
	defer resp.Body.Close()

	return nil
}

func (b *Basestation) GetState() (response *EventStreamResponse, err error) {

	payload := EventStreamPayload{
		Action:          "get",
		Resource:        "basestation",
		PublishResponse: false,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to get basestation state")
}

func (b *Basestation) GetAssociatedCamerasState() (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "get",
		Resource:        "cameras",
		PublishResponse: false,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to get associated cameras state")
}

func (b *Basestation) GetRules() (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "get",
		Resource:        "rules",
		PublishResponse: false,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to get rules")
}

func (b *Basestation) GetCalendarMode() (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "get",
		Resource:        "schedule",
		PublishResponse: false,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to get schedule")
}

// SetCalendarMode toggles calendar mode.
// NOTE: The Arlo API seems to disable calendar mode when switching to other modes, if it's enabled.
// You should probably do the same, although, the UI reflects the switch from calendar mode to say armed mode without explicitly setting calendar mode to inactive.
func (b *Basestation) SetCalendarMode(active bool) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        "schedule",
		PublishResponse: true,
		Properties: BasestationScheduleProperties{
			Active: active,
		},
		From: fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:   b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to set schedule")
}

func (b *Basestation) GetModes() (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "get",
		Resource:        "modes",
		PublishResponse: false,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to get modes")
}

func (b *Basestation) SetCustomMode(mode string) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        "modes",
		PublishResponse: true,
		Properties: BasestationModeProperties{
			Active: mode,
		},
		From: fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:   b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to set mode")
}

func (b *Basestation) DeleteMode(mode string) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "delete",
		Resource:        fmt.Sprintf("modes/%s", mode),
		PublishResponse: true,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to set mode")
}

func (b *Basestation) Arm() (response *EventStreamResponse, err error) {
	return b.SetCustomMode("mode1")
}

func (b *Basestation) Disarm() (response *EventStreamResponse, err error) {
	return b.SetCustomMode("mode0")
}

func (b *Basestation) SirenOn() (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        "siren",
		PublishResponse: true,
		Properties: SirenProperties{
			SirenState: "on",
			Duration:   300,
			Volume:     8,
			Pattern:    "alarm",
		},
		From: fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:   b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to get modes")
}

func (b *Basestation) SirenOff() (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        "siren",
		PublishResponse: true,
		Properties: SirenProperties{
			SirenState: "off",
			Duration:   300,
			Volume:     8,
			Pattern:    "alarm",
		},
		From: fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:   b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to get modes")
}
