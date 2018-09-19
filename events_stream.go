package arlo

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

type EventStream struct {
	SSEClient     *sse.Client
	Subscriptions map[string]chan *EventStreamResponse
	Events        chan *sse.Event
	ErrorChan     chan error
	Registered    bool
	Connected     bool
	Verbose       bool

	sync.Mutex
}

func NewEventStream(url string, client *http.Client) *EventStream {

	SSEClient := sse.NewClient(url)
	SSEClient.Connection = client

	return &EventStream{
		SSEClient:     SSEClient,
		Events:        make(chan *sse.Event),
		Subscriptions: make(map[string]chan *EventStreamResponse),
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
			/*
				fmt.Println("Got event message.")
				fmt.Printf("EVENT: %s\n", event.Event)
				fmt.Printf("DATA: %s\n", event.Data)
			*/

			if event.Data != nil {
				notifyResponse := &EventStreamResponse{}
				b := bytes.NewBuffer(event.Data)
				err := json.NewDecoder(b).Decode(notifyResponse)
				if err != nil {
					e.ErrorChan <- FAILED_TO_DECODE_JSON
					break
				}

				if notifyResponse.Status == "connected" {
					e.Connected = true
				} else if notifyResponse.Status == "disconnected" {
					e.Connected = false
				} else {
					if subscriber, ok := e.Subscriptions[notifyResponse.TransId]; ok {
						e.Lock()
						subscriber <- notifyResponse
						close(subscriber)
						delete(e.Subscriptions, notifyResponse.TransId)
						e.Unlock()
					}
				}
			}
		}
	}()
}
