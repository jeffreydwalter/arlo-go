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

// URL is part of the Status message fragment returned by most calls to the Arlo API.
// URL is only populated when Success is false.
type Data struct {
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Status is the message fragment returned from most http calls to the Arlo API.
type Status struct {
	Data    `json:"URL,omitempty"`
	Success bool `json:"success"`
}

// LoginResponse is an intermediate struct used when parsing data from the Login() call.
type LoginResponse struct {
	Data Account
	Status
}

type SessionResponse struct {
	Data Session
	Status
}

type UserProfileResponse struct {
	Data UserProfile
	Status
}

// DeviceResponse is an intermediate struct used when parsing data from the GetDevices() call.
type DeviceResponse struct {
	Data Devices
	Status
}

// LibraryMetaDataResponse is an intermediate struct used when parsing data from the GetLibraryMetaData() call.
type LibraryMetaDataResponse struct {
	Data LibraryMetaData
	Status
}

type LibraryResponse struct {
	Data Library
	Status
}

type CvrPlaylistResponse struct {
	Data CvrPlaylist
	Status
}

type Stream struct {
	URL string `json:"url"`
}

type StreamResponse struct {
	Data Stream
	Status
}

type RecordingResponse struct {
	Data Stream
	Status
}

type EventStreamResponse struct {
	EventStreamPayload
	Status string `json:"status,omitempty"`
}
