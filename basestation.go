package arlo

import (
	"fmt"

	"github.com/pkg/errors"
)

type BaseStationMetadata struct {
	InterfaceVersion         int             `json:"interfaceVersion"`
	ApiVersion               int             `json:"apiVersion"`
	State                    string          `json:"state"`
	SwVersion                string          `json:"swVersion"`
	HwVersion                string          `json:"hwVersion"`
	ModelId                  string          `json:"modelId"`
	Capabilities             []string        `json:"capabilities"`
	McsEnabled               bool            `json:"mcsEnabled"`
	AutoUpdateEnabled        bool            `json:"autoUpdateEnabled"`
	TimeZone                 string          `json:"timeZone"`
	OlsonTimeZone            string          `json:"olsonTimeZone"`
	UploadBandwidthSaturated bool            `json:"uploadBandwidthSaturated"`
	AntiFlicker              map[string]int  `json:"antiFlicker"`
	LowBatteryAlert          map[string]bool `json:"lowBatteryAlert"`
	LowSignalAlert           map[string]bool `json:"lowSignalAlert"`
	Claimed                  bool            `json:"claimed"`
	TimeSyncState            string          `json:"timeSyncState"`
	Connectivity             []struct {
		Type      string `json:"type"`
		Connected bool   `json:"connected"`
	} `json:"connectivity"`
}

// A Basestation is a Device that's not type "camera" (basestation, arloq, arloqs, etc.).
// This type is here just for semantics. Some methods explicitly require a device of a certain type.
type Basestation struct {
	Device
	eventStream *EventStream
}

// Basestations is an array of Basestation objects.
type Basestations []Basestation

func (b *Basestation) Subscribe() error {
	b.eventStream = NewEventStream(BaseUrl+fmt.Sprintf(SubscribeUri, b.arlo.Account.Token), b.arlo.client.HttpClient)
	connected := b.eventStream.Listen()

outoffor:
	for {
		// TODO: Need to add a timeout here.
		// We blocking here because we can't really do anything with the event stream until we're connected.
		// Once we have confirmation that we're connected to the event stream, we will "subscribe" to events.
		select {
		case b.eventStream.Connected = <-connected:
			if b.eventStream.Connected {
				break outoffor
			} else {
				// TODO: What do we do if Connected is false? Probably need retry logic here.
				break
			}
		case <-b.eventStream.Close:
			return errors.New("failed to subscribe to the event stream")
		}
	}

	// This is a crude (temporary?) way to monitor the connection. It's late and I'm tired, so this will probably go away.
	go func() {
	outoffor:
		for {
			select {
			case b.eventStream.Connected = <-connected:
				// TODO: What do we do if Connected is false? Probably need retry logic here.
				break outoffor
			case <-b.eventStream.Close:
				// TODO: Figure out what to do here if the eventStream is closed. (Panic?)
				return
			}
		}
	}()

	payload := Payload{
		Action:          "set",
		Resource:        fmt.Sprintf("subscriptions/%s_%s", b.UserId, TransIdPrefix),
		PublishResponse: false,
		Properties:      map[string][1]string{"devices": {b.DeviceId}},
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	if _, err := b.makeEventStreamRequest(payload, "failed to subscribe to the event stream"); err != nil {
		return err
	}

	return nil
}

func (b *Basestation) Unsubscribe() error {
	// TODO: Close channel to stop EventStream.
	//return errors.New("not implemented")
	if b.eventStream != nil {
		close(b.eventStream.Close)
	}
	return nil
}

func (b *Basestation) IsConnected() error {
	if !b.eventStream.Connected {
		return errors.New("basestation not connected to event stream")
	}
	return nil
}

func (b *Basestation) GetState() (*EventStreamResponse, error) {

	payload := Payload{
		Action:          "get",
		Resource:        "basestation",
		PublishResponse: false,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to get basestation state")
}

func (b *Basestation) GetAssociatedCamerasState() (*EventStreamResponse, error) {
	payload := Payload{
		Action:          "get",
		Resource:        "cameras",
		PublishResponse: false,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	return b.makeEventStreamRequest(payload, "failed to get associated cameras state")
}

func (b *Basestation) makeEventStreamRequest(payload Payload, msg string) (*EventStreamResponse, error) {
	transId := genTransId()
	payload.TransId = transId

	if err := b.IsConnected(); err != nil {
		return nil, errors.WithMessage(err, msg)
	}

	b.eventStream.Subscriptions[transId] = make(chan *EventStreamResponse)
	defer close(b.eventStream.Subscriptions[transId])

	resp, err := b.arlo.post(fmt.Sprintf(NotifyUri, b.DeviceId), b.XCloudId, payload, nil)
	if err := checkRequest(*resp, err, msg); err != nil {
		return nil, err
	}

	select {
	case eventStreamResponse := <-b.eventStream.Subscriptions[transId]:
		return eventStreamResponse, nil
	case err = <-b.eventStream.Error:
		return nil, errors.Wrap(err, "failed to get basestation")
	case <-b.eventStream.Close:
		return nil, errors.New("event stream was closed before response was read")
	}
}
