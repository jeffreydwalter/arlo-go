package arloclient

// Credentials is the login credential data.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Duration holds two dates used when you need to specify a date range in the format "20060102".
type Duration struct {
	DateFrom string `json:"dateFrom""`
	DateTo   string `json:"dateTo"`
}

// PasswordPair is used when updating the account password.
type PasswordPair struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// FullName is used when updating the account username.
type FullName struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// Account is the account data.
type Account struct {
	UserId        string `json:"userId"`
	Email         string `json:"email"`
	Token         string `json:"token"`
	PaymentId     string `json:"paymentId"`
	Authenticated uint32 `json:"authenticated"`
	AccountStatus string `json:"accountStatus"`
	SerialNumber  string `json:"serialNumber"`
	CountryCode   string `json:"countryCode"`
	TocUpdate     bool   `json:"tocUpdate"`
	PolicyUpdate  bool   `json:"policyUpdate"`
	ValidEmail    bool   `json:"validEmail"`
	Arlo          bool   `json:"arlo"`
	DateCreated   uint64 `json:"dateCreated"`
}

// Owner is the owner of a Device data.
type Owner struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	OwnerId   string `json:"ownerId"`
}

// Properties is the Device properties data.
type Properties struct {
	ModelId       string `json:"modelId"`
	OlsonTimeZone string `json:"olsonTimeZone"`
	HwVersion     string `json:"hwVersion"`
}

// Device is the device data.
type Device struct {
	DeviceType         string     `json:"deviceType"`
	XCloudId           string     `json:"xCloudId"`
	DisplayOrder       uint8      `json:"displayOrder"`
	State              string     `json:"state"`
	ModelId            string     `json:"modelId"`
	InterfaceVersion   string     `json:"interfaceVersion"`
	UserId             string     `json:"userId"`
	DeviceName         string     `json:"deviceName"`
	FirmwareVersion    string     `json:"firmwareVersion"`
	MediaObjectCount   uint8      `json:"mediaObjectCount"`
	DateCreated        uint64     `json:"dateCreated"`
	Owner              Owner      `json:"owner"`
	Properties         Properties `json:"properties"`
	UniqueId           string     `json:"uniqueId"`
	LastModified       float64    `json:"lastModified"`
	UserRole           string     `json:"userRole"`
	InterfaceSchemaVer string     `json:"interfaceSchemaVer"`
	DeviceId           string     `json:"deviceId"`
}

// Devices is an array of Device objects.
type Devices []Device

// LibraryMetaData is the library meta data.
type LibraryMetaData struct {
	// TODO: Fill this out.
}

// UserProfile is the user profile data.
type UserProfile struct {
	// TODO: Fill this out.
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

// UserProfile is an intermediate struct used when parsing data from the UpdateProfile() call.
type UserProfileResponse struct {
	Data    UserProfile
	Success bool `json:"success"`
}
