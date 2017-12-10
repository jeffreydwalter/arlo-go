package arlo

/*
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
*/

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

type Favorite struct {
	NonFavorite uint8 `json:"nonFavorite"`
	Favorite    uint8 `json:"Favorite"`
}

/*
type Device struct {
	DeviceType         string     `json:"deviceType"`
	XCloudId           string     `json:"xCloudId"`
	DisplayOrder       uint8      `json:"displayOrder"`
	State              string     `json:"state"`
	ModelId            string     `json:"modelId"`
	InterfaceVersion   string     `json:"interfaceVersion"`
	ParentId           string     `json:"parentId"`
	UserId             string     `json:"userId"`
	DeviceName         string     `json:"deviceName"`
	FirmwareVersion    string     `json:"firmwareVersion"`
	MediaObjectCount   uint8      `json:"mediaObjectCount"`
	DateCreated        float64    `json:"dateCreated"`
	Owner              Owner      `json:"owner"`
	Properties         Properties `json:"properties"`
	UniqueId           string     `json:"uniqueId"`
	LastModified       float64    `json:"lastModified"`
	UserRole           string     `json:"userRole"`
	InterfaceSchemaVer string     `json:"interfaceSchemaVer"`
	DeviceId           string     `json:"deviceId"`
}
*/

type StreamUrl struct {
	Url string `json:"url"`
}

type NotificationProperties struct {
	ActivityState string `json:"activityState"`
	CameraId      string `json:"cameraId"`
}

type Notification struct {
	To              string                 `json:"to"`
	From            string                 `json:"from"`
	Resource        string                 `json:"resource"`
	Action          string                 `json:"action"`
	PublishResponse bool                   `json:"publishResourcec"`
	TransId         string                 `json:"transId"`
	Properties      NotificationProperties `json:"properties"`
}
