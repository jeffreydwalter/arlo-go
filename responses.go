package arlo

// LoginResponse is an intermediate struct used when parsing data from the Login() call.
type LoginResponse struct {
	Data Account
	Error
}

// DeviceResponse is an intermediate struct used when parsing data from the GetDevices() call.
type DeviceResponse struct {
	Data Devices
	Error
}

// LibraryMetaDataResponse is an intermediate struct used when parsing data from the GetLibraryMetaData() call.
type LibraryMetaDataResponse struct {
	Data LibraryMetaData
	Error
}

type LibraryResponse struct {
	Data Library
	Error
}

type StreamResponse struct {
	Data StreamUrl
	Error
}

type RecordingResponse struct {
	Data StreamUrl
	Error
}

type EventStreamResponse struct {
	Action     string      `json:"action,omitempty"`
	Resource   string      `json:"resource,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
	TransId    string      `json:"transId"`
	From       string      `json:"from"`
	To         string      `json:"to"`
	Status     string      `json:"status"`
}
