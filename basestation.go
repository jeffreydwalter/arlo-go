package arlo

import (
	"fmt"
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
	b.eventStream.Listen()

	body := Payload{
		Action:          "set",
		Resource:        fmt.Sprintf("subscriptions/%s_%s", b.UserId, TransIdPrefix),
		PublishResponse: false,
		Properties:      map[string][1]string{"devices": {b.DeviceId}},
		TransId:         genTransId(),
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}
	resp, err := b.arlo.post(fmt.Sprintf(NotifyUri, b.DeviceId), b.XCloudId, body, nil)
	return checkRequest(*resp, err, "failed to subscribe to the event stream")
}

func (b *Basestation) Unsubscribe() error {
	// TODO: Close channel to stop EventStream.
	//return errors.New("not implemented")
	return nil
}

func (b *Basestation) GetState() (*EventStreamResponse, error) {
	transId := genTransId()

	b.eventStream.Subscriptions[transId] = make(chan *EventStreamResponse)

	body := Payload{
		Action:          "get",
		Resource:        "basestation",
		PublishResponse: false,
		TransId:         transId,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	resp, err := b.arlo.post(fmt.Sprintf(NotifyUri, b.DeviceId), b.XCloudId, body, nil)
	if err := checkRequest(*resp, err, "failed to get basestation state"); err != nil {
		return nil, err
	}

	return <-b.eventStream.Subscriptions[transId], nil
}

func (b *Basestation) GetAssociatedCamerasState() (*EventStreamResponse, error) {
	transId := genTransId()

	b.eventStream.Subscriptions[transId] = make(chan *EventStreamResponse)

	body := Payload{
		Action:          "get",
		Resource:        "cameras",
		PublishResponse: false,
		TransId:         transId,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	resp, err := b.arlo.post(fmt.Sprintf(NotifyUri, b.DeviceId), b.XCloudId, body, nil)
	if err := checkRequest(*resp, err, "failed to get camera state"); err != nil {
		return nil, err
	}

	return <-b.eventStream.Subscriptions[transId], nil
}
