package arlo_golang

import (
	"encoding/json"
	"fmt"
	"time"

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
}

// Basestations is an array of Basestation objects.
type Basestations []Basestation

func (b *Basestation) connect(a *Arlo) {
	b.eventStream = NewEventStream(BaseUrl+fmt.Sprintf(SubscribeUri, a.Account.Token), util.HeaderToMap(*a.client.BaseHttpHeader))
	b.eventStream.Listen()
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
func (a *Arlo) GetBasestationState(b Basestation) (*NotifyResponse, error) {

	transId := GenTransId()

	body := NotifyPayload{
		Action:          "get",
		Resource:        "basestation",
		PublishResponse: false,
		Properties:      map[string]string{},
		TransId:         transId,
		From:            fmt.Sprintf("%s_%s", b.UserId, TransIdPrefix),
		To:              b.DeviceId,
	}

	b.eventStream.Subscriptions[transId] = new(Subscriber)

	for b.eventStream.Connected == false {
		fmt.Println("Not connected yet.")
		time.Sleep(1000 * time.Millisecond)
	}
	fmt.Println("Connected now.")

	resp, err := a.client.Post(fmt.Sprintf(NotifyUri, b.DeviceId), body, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to start stream")
	}

	ep := &NotifyResponse{}
	err = json.NewDecoder(resp.Body).Decode(ep)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to decode body")
	}

	for {
		fmt.Println("Subscribing to the eventstream.")
		select {
		case notifyResponse := <-*b.eventStream.Subscriptions[transId]:
			fmt.Println("Recieved a response from the subscription.")
			return &notifyResponse, nil
		}
	}
}
