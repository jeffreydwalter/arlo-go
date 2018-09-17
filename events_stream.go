package arlo_golang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/pkg/errors"

	"github.com/r3labs/sse"
)

var FAILED_TO_PUBLISH = errors.New("Failed to publish")

var FAILED_TO_DECODE_JSON = errors.New("Failed to decode JSON")

var FAILED_TO_SUBSCRIBE = errors.New("Failed to subscribe to SSEClient")

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

func NewEventStream(url string, headers map[string]string) *EventStream {

	SSEClient := sse.NewClient(url)
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
		err := e.SSEClient.SubscribeChan("", e.Events)
		if err != nil {
			fmt.Println(FAILED_TO_SUBSCRIBE)
			e.ErrorChan <- FAILED_TO_SUBSCRIBE
		}
	}()

	for event := range e.Events {
		fmt.Println("Got event message here.")
		fmt.Printf("EVENT: %s\n", event.Event)
		fmt.Printf("DATA: %s\n", event.Data)

		if event.Data != nil {
			notifyResponse := &NotifyResponse{}
			b := bytes.NewBuffer(event.Data)
			err := json.NewDecoder(b).Decode(notifyResponse)
			if err != nil {
				e.ErrorChan <- errors.WithMessage(err, "failed to decode JSON")
				break
			}

			if notifyResponse.Status == "connected" {
				e.Connected = true
				fmt.Println("Connected.")
				break
			}
		}
	}

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

				if notifyResponse.Status == "connected" {
					fmt.Println("Connected.")
					e.Connected = true
				} else if notifyResponse.Status == "disconnected" {
					fmt.Println("Disconnected.")
					e.Connected = false
				} else {
					fmt.Printf("Message for transId: %s", notifyResponse.TransId)
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
			}
		}
	}()
	/*
		go func() {

			fmt.Println("go func to recieve a subscription.")
			for {
				fmt.Println("go func for loop to recieve a subscription.")
				select {
				case s := <-e.Subscriptions:
					if resp, ok := e.Responses[s.transId]; ok {
						fmt.Println("Recieved a subscription, sending response.")
						s.ResponseChan <- resp
						e.Lock()
						delete(e.Responses, s.transId)
						e.Unlock()
					} else {
						fmt.Println("Recieved a subscription error, sending error response.")
						e.ErrorChan <- FAILED_TO_PUBLISH
						break
					}
				}
			}
		}()
	*/
}

func (e *EventStream) verbose(params ...interface{}) {
	if e.Verbose {
		log.Println(params...)
	}
}
