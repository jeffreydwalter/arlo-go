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
	DateCreated   int64  `json:"dateCreated"`
}

// Friend is the account data for non-primary account holders designated as friends.
type Friend struct {
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Devices      DeviceOrder `json:"devices"`
	LastModified int64       `json:"lastModified"`
	AdminUser    bool        `json:"adminUser"`
	Email        string      `json:"email"`
	Id           string      `json:"id"`
}

// Connectivity is part of the Device data.
type Connectivity struct {
	ActiveNetwork  string `json:"activeNetwork,omitempty"`
	APN            string `json:"apn,omitempty"`
	CarrierFw      string `json:"carrierFw,omitempty"`
	Connected      bool   `json:"connected,omitempty"`
	FWVersion      string `json:"fwVersion,omitempty"`
	ICCID          string `json:"iccid,omitempty"`
	IMEI           string `json:"imei,omitempty"`
	MEPStatus      string `json:"mepStatus,omitempty"`
	MSISDN         string `json:"msisdn,omitempty"`
	NetworkMode    string `json:"networkMode,omitempty"`
	NetworkName    string `json:"networkName,omitempty"`
	RFBand         int    `json:"rfBand,omitempty"`
	Roaming        bool   `json:"roaming"`
	RoamingAllowed bool   `json:"roamingAllowed"`
	SignalStrength string `json:"signalStrength,omitempty"`
	Type           string `json:"type,omitempty"`
	WWANIPAddr     string `json:"wwanIpAddr,omitempty"`
}

type BaseStationMetadata struct {
	InterfaceVersion         int             `json:"interfaceVersion"`
	ApiVersion               int             `json:"apiVersion"`
	State                    string          `json:"state"`
	SwVersion                string          `json:"swVersion"`
	HwVersion                string          `json:"hwVersion"`
	ModelId                  string          `json:"modelId"`
	Capabilities             []string        `json:"capabilities"`
	McsEnabled               bool            `json:"mcsEnabled"`
	AutoUpdateEnabled        bool            `json:"autoUpdateEnabled"`
	TimeZone                 string          `json:"timeZone"`
	OlsonTimeZone            string          `json:"olsonTimeZone"`
	UploadBandwidthSaturated bool            `json:"uploadBandwidthSaturated"`
	AntiFlicker              map[string]int  `json:"antiFlicker"`
	LowBatteryAlert          map[string]bool `json:"lowBatteryAlert"`
	LowSignalAlert           map[string]bool `json:"lowSignalAlert"`
	Claimed                  bool            `json:"claimed"`
	TimeSyncState            string          `json:"timeSyncState"`
	Connectivity             Connectivity    `json:"connectivity"`
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

type BaseDetectionProperties struct {
	Armed       bool     `json:"armed"`
	Sensitivity int      `json:"sensitivity"`
	Zones       []string `json:"zones,omitempty"`
}

// MotionDetectionProperties is the Properties struct for the EventStreamPayload type.
type MotionDetectionProperties struct {
	BaseDetectionProperties `json:"motionDetection"`
}

// AudioDetectionProperties is the Properties struct for the EventStreamPayload type.
type AudioDetectionProperties struct {
	BaseDetectionProperties `json:"audioDetection"`
}

type EmailNotification struct {
	Enabled          bool     `json:"enabled"`
	EmailList        []string `json:"emailList"`
	PushNotification bool     `json:"pushNotification"`
}

type PlayTrackProperties struct {
	TrackId  string `json:"trackId"`
	Position int    `json:"position"`
}

type BaseLoopbackModeProperties struct {
	LoopbackMode string `json:"loopbackMode"`
}

type LoopbackModeProperties struct {
	Config BaseLoopbackModeProperties `json:"config"`
}

type BaseSleepTimerProperties struct {
	SleepTime    int64 `json:"sleepTime"`
	SleepTimeRel int   `json:"sleepTimeRel"`
}

type SleepTimerProperties struct {
	Config BaseSleepTimerProperties `json:"config"`
}
type BaseEventActionProperties struct {
	ActionType        string `json:"actionType"`
	StopType          string `json:"stopType"`
	Timeout           int    `json:"timeout"`
	EmailNotification `json:"emailNotification"`
}

type EventActionProperties struct {
	BaseEventActionProperties `json:"eventAction"`
}

type BaseShuffleProperties struct {
	ShuffleActive bool `json:"shuffleActive"`
}

type ShuffleProperties struct {
	Config BaseShuffleProperties `json:"config"`
}

type VolumeProperties struct {
	Mute   bool `json:"mute"`
	Volume int  `json:"volume,omitempty"`
}

type SpeakerProperties struct {
	Speaker VolumeProperties `json:"speaker"`
}

type NightLightRGBProperties struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}

type BaseNightLightProperties struct {
	Brightness   int                     `json:"brightness,omitempty"`
	Enabled      bool                    `json:"enabled"`
	Mode         string                  `json:"mode,omitempty"`
	RGB          NightLightRGBProperties `json:"mode,omitempty"`
	SleepTime    int64                   `json:"sleepTime,omitempty"`
	SleepTimeRel int                     `json:"sleepTimeRel,omitempty"`
}

type NightLightProperties struct {
	NightLight BaseNightLightProperties `json:"nightLight"`
}

type SirenProperties struct {
	SirenState string `json:"sirenState"`
	Duration   int    `json:"duration"`
	Volume     int    `json:"volume"`
	Pattern    string `json:"pattern"`
}

type BasestationModeProperties struct {
	Active string `json:"active"`
}

type BasestationScheduleProperties struct {
	Active bool `json:"active"`
}

type CameraProperties struct {
	PrivacyActive bool `json:"privacyActive"`
	Brightness    int  `json:"brightness,omitempty"`
}

// EventStreamPayload is the message that will be sent to the arlo servers via the /notify API.
type EventStreamPayload struct {
	Action          string      `json:"action,omitempty"`
	Resource        string      `json:"resource,omitempty"`
	PublishResponse bool        `json:"publishResponse"`
	Properties      interface{} `json:"properties,omitempty"`
	TransId         string      `json:"transId"`
	From            string      `json:"from"`
	To              string      `json:"to"`
}

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
