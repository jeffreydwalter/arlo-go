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

// LoginResponse is an intermediate struct used when parsing data from the Login() call.
type LoginResponse struct {
	Data Account
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

type StreamResponse struct {
	URL string `json:"url"`
	Status
}

type RecordingResponse struct {
	URL string `json:"url"`
	Status
}

type EventStreamResponse struct {
	EventStreamPayload
	Status string `json:"status,omitempty"`
}
