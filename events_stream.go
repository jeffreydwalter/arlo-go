package arlo_golang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/pkg/errors"

	"github.com/r3labs/sse"
)

var (
	FAILED_TO_PUBLISH     = errors.New("Failed to publish")
	FAILED_TO_DECODE_JSON = errors.New("Failed to decode JSON")
	FAILED_TO_SUBSCRIBE   = errors.New("Failed to subscribe to SSEClient")
)

type Subscriber chan NotifyResponse

type EventStream struct {
	Registered    bool
	Connected     bool
	SSEClient     *sse.Client
	Events        chan *sse.Event
	Subscriptions map[string]*Subscriber
	ErrorChan     chan error
	Responses     map[string]NotifyResponse
	Verbose       bool

	sync.Mutex
}

func NewEventStream(url string, client *http.Client, headers map[string]string) *EventStream {

	SSEClient := sse.NewClient(url)
	SSEClient.Connection = client
	SSEClient.Headers = headers

	return &EventStream{
		SSEClient:     SSEClient,
		Events:        make(chan *sse.Event),
		Subscriptions: map[string]*Subscriber{},
		ErrorChan:     make(chan error, 1),
	}
}

func (e *EventStream) Listen() {

	go func() {
		err := e.SSEClient.SubscribeChanRaw(e.Events)
		if err != nil {
			fmt.Println(FAILED_TO_SUBSCRIBE)
			e.ErrorChan <- FAILED_TO_SUBSCRIBE
		}
	}()

	go func() {
		for event := range e.Events {
			fmt.Println("Got event message.")
			fmt.Printf("EVENT: %X\n", event.Event)
			fmt.Printf("DATA: %X\n", event.Data)

			if event.Data != nil {
				notifyResponse := &NotifyResponse{}
				b := bytes.NewBuffer(event.Data)
				err := json.NewDecoder(b).Decode(notifyResponse)
				if err != nil {
					e.ErrorChan <- FAILED_TO_DECODE_JSON
					break
				}

				fmt.Printf("%s\n", notifyResponse)
				if notifyResponse.Status == "connected" {
					fmt.Println("Connected.")
					e.Connected = true
				} else if notifyResponse.Status == "disconnected" {
					fmt.Println("Disconnected.")
					e.Connected = false
				} else {
					fmt.Printf("Message for transId: %s\n", notifyResponse.TransId)
					if subscriber, ok := e.Subscriptions[notifyResponse.TransId]; ok {
						e.Lock()
						*subscriber <- *notifyResponse
						close(*subscriber)
						delete(e.Subscriptions, notifyResponse.TransId)
						e.Unlock()
					} else {
						// Throw away the message.
						fmt.Println("Throwing away message.")
					}
				}
			} else {
				fmt.Printf("Event data was nil.\n")
			}
		}
	}()
}
