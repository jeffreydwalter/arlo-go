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
	FAILED_TO_PUBLISH     = errors.New("failed to publish")
	FAILED_TO_DECODE_JSON = errors.New("failed to decode json")
	FAILED_TO_SUBSCRIBE   = errors.New("failed to subscribe to seeclient")
)

type EventStream struct {
	SSEClient     *sse.Client
	Subscriptions map[string]chan *EventStreamResponse
	Events        chan *sse.Event
	Error         chan error
	Close         chan interface{}
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
		Error:         make(chan error),
		Close:         make(chan interface{}),
	}
}

func (e *EventStream) Listen() (connected chan bool) {

	connected = make(chan bool)

	go func() {
		err := e.SSEClient.SubscribeChanRaw(e.Events)
		if err != nil {
			fmt.Println(FAILED_TO_SUBSCRIBE)
			e.Error <- FAILED_TO_SUBSCRIBE
		}

		for {
			select {
			case event := <-e.Events:
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
						e.Error <- FAILED_TO_DECODE_JSON
						break
					}

					if notifyResponse.Status == "connected" {
						connected <- true
					} else if notifyResponse.Status == "disconnected" {
						connected <- false
					} else {
						if subscriber, ok := e.Subscriptions[notifyResponse.TransId]; ok {
							e.Lock()
							subscriber <- notifyResponse
							e.Unlock()
						}
					}
				}
			case <-e.Close:
				connected <- false
				return
			}
		}
	}()

	return connected
}
