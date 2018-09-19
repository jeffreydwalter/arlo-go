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

// Account is the account data.
type Account struct {
	UserId        string  `json:"userId"`
	Email         string  `json:"email"`
	Token         string  `json:"token"`
	PaymentId     string  `json:"paymentId"`
	Authenticated uint32  `json:"authenticated"`
	AccountStatus string  `json:"accountStatus"`
	SerialNumber  string  `json:"serialNumber"`
	CountryCode   string  `json:"countryCode"`
	TocUpdate     bool    `json:"tocUpdate"`
	PolicyUpdate  bool    `json:"policyUpdate"`
	ValidEmail    bool    `json:"validEmail"`
	Arlo          bool    `json:"arlo"`
	DateCreated   float64 `json:"dateCreated"`
}

type Friend struct {
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Devices      DeviceOrder `json:"devices"`
	LastModified float64     `json:"lastModified"`
	AdminUser    bool        `json:"adminUser"`
	Email        string      `json:"email"`
	Id           string      `json:"id"`
}

// Owner is part of the Device data.
type Connectivity struct {
	ActiveNetwork  string `json:"activeNetwork"`
	APN            string `json:"apn"`
	CarrierFw      string `json:"carrierFw"`
	Connected      bool   `json:"connected"`
	FWVersion      string `json:"fwVersion"`
	ICCID          string `json:"iccid"`
	IMEI           string `json:"imei"`
	MEPStatus      string `json:"mepStatus"`
	MSISDN         string `json:"msisdn"`
	NetworkMode    string `json:"networkMode"`
	NetworkName    string `json:"networkName"`
	RFBand         int    `json:"rfBand"`
	Roaming        bool   `json:"roaming"`
	RoamingAllowed bool   `json:"roamingAllowed"`
	SignalStrength string `json:"signalStrength"`
	Type           string `json:"type"`
	WWANIPAddr     string `json:"wwanIpAddr"`
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

type Favorite struct {
	NonFavorite uint8 `json:"nonFavorite"`
	Favorite    uint8 `json:"Favorite"`
}

type StreamUrl struct {
	Url string `json:"url"`
}

// Payload represents the message that will be sent to the arlo servers via the Notify API.
type Payload struct {
	Action          string      `json:"action,omitempty"`
	Resource        string      `json:"resource,omitempty"`
	PublishResponse bool        `json:"publishResponse"`
	Properties      interface{} `json:"properties,omitempty"`
	TransId         string      `json:"transId"`
	From            string      `json:"from"`
	To              string      `json:"to"`
}

type Data struct {
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Error   string `json:"error,omitempty"`
}

// map[data:map[message:The device does not exist. reason:No such device. error:2217] success:false]
type Error struct {
	Data    `json:"Data,omitempty"`
	Success bool `json:"success"`
}
