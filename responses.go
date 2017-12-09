package arloclient

// UpdateResponse is an intermediate struct used when parsing data from the UpdateProfile() call.
type Status struct {
	Success bool `json:"success"`
}

// LoginResponse is an intermediate struct used when parsing data from the Login() call.
type LoginResponse struct {
	Data    Account
	Success bool `json:"success"`
}

// DeviceResponse is an intermediate struct used when parsing data from the GetDevices() call.
type DeviceResponse struct {
	Data    Devices
	Success bool `json:"success"`
}

// LibraryMetaDataResponse is an intermediate struct used when parsing data from the GetLibraryMetaData() call.
type LibraryMetaDataResponse struct {
	Data    LibraryMetaData
	Success bool `json:"success"`
}

type LibraryResponse struct {
	Data    Library
	Success bool `json:"success"`
}
