/*
 * Copyright (c) 2018 Jeffrey Walter <jeffreydwalter@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

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
		e.SSEClient.OnDisconnect(func(c *sse.Client) {
			e.disconnect()
			// fmt.Printf("\n\n\n\nClIENT DISCONNECTED!!!!!\n\n\n\n")
		})
		err := e.SSEClient.SubscribeChanRaw(e.Events)
		if err != nil {
			e.Error <- FAILED_TO_SUBSCRIBE
		}

		for {
			select {
			case event := <-e.Events:
				//fmt.Println("Got event message.")
				/*
					fmt.Print(".")
					fmt.Printf("EVENT: %s\n", event.Event)
					fmt.Printf("DATA: %s\n", event.Data)
				*/

				if event != nil && event.Data != nil {
					notifyResponse := &EventStreamResponse{}
					b := bytes.NewBuffer(event.Data)
					err := json.NewDecoder(b).Decode(notifyResponse)
					if err != nil {
						e.Error <- FAILED_TO_DECODE_JSON
						break
					}

					// FIXME: This is a shitty way to handle this. It's potentially leaking a chan.
					if notifyResponse.Status == "connected" {
						connected <- true
					} else if notifyResponse.Status == "disconnected" {
						e.disconnect()
                    			} else if notifyResponse.Action == "logout" {
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
