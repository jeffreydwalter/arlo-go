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
	Data StreamUrl
	Status
}

type RecordingResponse struct {
	Data StreamUrl
	Status
}

type EventStreamResponse struct {
	EventStreamPayload
	Status string `json:"status,omitempty"`
}
