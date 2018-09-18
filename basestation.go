package arlo_golang

import (
	"fmt"

	"github.com/jeffreydwalter/arlo-golang/internal/util"
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
	arlo        *Arlo
}

// Basestations is an array of Basestation objects.
type Basestations []Basestation

func (b *Basestation) Subscribe() (*Status, error) {
	b.eventStream = NewEventStream(BaseUrl+fmt.Sprintf(SubscribeUri, b.arlo.Account.Token), b.arlo.client.HttpClient, util.HeaderToMap(*b.arlo.client.BaseHttpHeader))
	b.eventStream.Listen()

	transId := GenTransId()

	body := NotifyPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("subscriptions/%s_%s", b.UserId, "web"),
		PublishResponse: false,
		Properties:      map[string][]string{"devices": []string{b.DeviceId}},
		TransId:         transId,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	resp, err := b.arlo.client.Post(fmt.Sprintf(NotifyUri, b.DeviceId), body, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to subscribe to the event stream")
	}

	var status Status
	if err := resp.Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}

/*
 This is an example of the json you would pass in the body to UpdateFriends():
{
  "firstName":"Some",
  "lastName":"Body",
  "devices":{
    "XXXXXXXXXXXXX":"Camera 1",
    "XXXXXXXXXXXXX":"Camera 2 ",
    "XXXXXXXXXXXXX":"Camera 3"
  },
  "lastModified":1463977440911,
  "adminUser":true,
  "email":"user@example.com",
  "id":"XXX-XXXXXXX"
}
*/
func (b *Basestation) GetState() (*NotifyResponse, error) {

	transId := GenTransId()

	body := NotifyPayload{
		Action:          "get",
		Resource:        "basestation",
		PublishResponse: false,
		TransId:         transId,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	//fmt.Printf("BODY: %+v\n", body)
	//fmt.Printf("HEADERS: %+v\n", a.client.BaseHttpHeader)

	fmt.Println("Subscribing to the eventstream.")
	b.eventStream.Subscriptions[transId] = new(Subscriber)

	resp, err := b.arlo.client.Post(fmt.Sprintf(NotifyUri, b.DeviceId), body, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get basestation state")
	}

	var status Status
	if err := resp.Decode(&status); err != nil {
		return nil, err
	}

	if !status.Success {
		return nil, errors.New("failed to get basestation status")
	}

	notifyResponse := <-*b.eventStream.Subscriptions[transId]
	return &notifyResponse, nil
}
