package arlo

import (
	"bytes"
	"encoding/json"
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

type subscriber chan *EventStreamResponse

type subscribers map[string]subscriber

type subscriptions struct {
	subscribers
	rwmutex sync.RWMutex
}

type eventStream struct {
	SSEClient    *sse.Client
	Events       chan *sse.Event
	Error        chan error
	Verbose      bool
	Disconnected chan interface{}
	once         *sync.Once

	subscriptions
}

func newEventStream(url string, client *http.Client) *eventStream {

	SSEClient := sse.NewClient(url)
	SSEClient.Connection = client

	return &eventStream{
		SSEClient:     SSEClient,
		Events:        make(chan *sse.Event),
		subscriptions: subscriptions{make(map[string]subscriber), sync.RWMutex{}},
		Error:         make(chan error),
		Disconnected:  make(chan interface{}),
		once:          new(sync.Once),
	}
}

func (e *eventStream) disconnect() {
	e.once.Do(func() {
		close(e.Disconnected)
	})
}

func (e *eventStream) listen() (connected chan bool) {

	connected = make(chan bool)

	go func() {
		err := e.SSEClient.SubscribeChanRaw(e.Events)
		if err != nil {
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
						e.disconnect()
					} else {
						e.subscriptions.rwmutex.RLock()
						subscriber, ok := e.subscribers[notifyResponse.TransId]
						e.subscriptions.rwmutex.RUnlock()
						if ok {
							subscriber <- notifyResponse
						}
					}
				}
			case <-e.Disconnected:
				connected <- false
				return
			}
		}
	}()

	return connected
}

func (s *subscriptions) unsubscribe(transId string) {
	s.rwmutex.Lock()
	defer s.rwmutex.Unlock()
	if _, ok := s.subscribers[transId]; ok {
		close(s.subscribers[transId])
		delete(s.subscribers, transId)
	}

}

func (s *subscriptions) subscribe(transId string, subscriber subscriber) {
	s.rwmutex.Lock()
	s.subscribers[transId] = subscriber
	s.rwmutex.Unlock()
}
