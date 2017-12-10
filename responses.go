package arlo

// UpdateResponse is an intermediate struct used when parsing data from the UpdateProfile() call.
type Status struct {
	Success bool `json:"success"`
}

// LoginResponse is an intermediate struct used when parsing data from the Login() call.
type LoginResponse struct {
	Data Account
	*Status
}

// DeviceResponse is an intermediate struct used when parsing data from the GetDevices() call.
type DeviceResponse struct {
	Data Devices
	*Status
}

// LibraryMetaDataResponse is an intermediate struct used when parsing data from the GetLibraryMetaData() call.
type LibraryMetaDataResponse struct {
	Data LibraryMetaData
	*Status
}

type LibraryResponse struct {
	Data Library
	*Status
}

type StartStreamResponse struct {
	Data StreamUrl
	*Status
}
