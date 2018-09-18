package arlo_golang

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

type StreamUrl struct {
	Url string `json:"url"`
}

// NotifyPayload represents the message that will be sent to the Arlo servers via the Notify API.
type NotifyPayload struct {
	Action          string      `json:"action,omitempty"`
	Resource        string      `json:"resource,omitempty"`
	PublishResponse bool        `json:"publishResponse"`
	Properties      interface{} `json:"properties,omitempty"`
	TransId         string      `json:"transId"`
	From            string      `json:"from"`
	To              string      `json:"to"`
}

type NotifyResponse struct {
	Action     string      `json:"action,omitempty"`
	Resource   string      `json:"resource,omitempty"`
	Properties interface{} `json:"properties,omitempty"`
	TransId    string      `json:"transId"`
	From       string      `json:"from"`
	To         string      `json:"to"`
	Status     string      `json:"status"`
}

/*
{"status":"connected"}

{"resource":"subscriptions/336-4764296_web","transId":"web!f94fbae4.46e6e!1520148142862","action":"is","from":"48935B7SA9847","to":"336-4764296_web","properties":{"devices":["48935B7SA9847"],"url":"https://vzweb05-prod.vz.netgear.com/hmsweb/publish/48935B7SA9847/336-4764296/c16ec5b2-f914-4140-aa5d-880feda292a0"}}

{"resource":"cameras/48B45974D8E54","properties":{"batteryLevel":15},"transId":"48935B7SA9847!cfa2b5ed!1520148143870","from":"48935B7SA9847","action":"is"}

{"resource":"basestation","properties":{"interfaceVersion":3,"apiVersion":1,"state":"idle","swVersion":"1.9.8.0_16666","hwVersion":"VMB3010r2","modelId":"VMB3010","capabilities":["gateway"],"mcsEnabled":true,"autoUpdateEnabled":true,"timeZone":"CST6CDT,M3.2.0,M11.1.0","olsonTimeZone":"America/Chicago","uploadBandwidthSaturated":false,"antiFlicker":{"mode":0,"autoDefault":1},"lowBatteryAlert":{"enabled":true},"lowSignalAlert":{"enabled":false},"claimed":true,"timeSyncState":"synchronized","connectivity":[{"type":"ethernet","connected":true}]},"action":"is","transId":"web!ffe75798.f6dca!1520148144127","to":"336-4764296_web","from":"48935B7SA9847"}


{"resource":"basestation","properties":{"interfaceVersion":3,"apiVersion":1,"state":"idle","swVersion":"1.9.8.0_16666","hwVersion":"VMB3010r2","modelId":"VMB3010","capabilities":["gateway"],"mcsEnabled":true,"autoUpdateEnabled":true,"timeZone":"CST6CDT,M3.2.0,M11.1.0","olsonTimeZone":"America/Chicago","uploadBandwidthSaturated":false,"antiFlicker":{"mode":0,"autoDefault":1},"lowBatteryAlert":{"enabled":true},"lowSignalAlert":{"enabled":false},"claimed":true,"timeSyncState":"synchronized","connectivity":[{"type":"ethernet","connected":true}]},"action":"is","transId":"web!ffe75798.f6dca!1520148144127","to":"336-4764296_web","from":"48935B7SA9847"}

{"resource":"cameras","properties":[{"interfaceVersion":3,"serialNumber":"48B45974D8E54","batteryLevel":15,"signalStrength":4,"brightness":0,"mirror":true,"flip":true,"powerSaveMode":3,"capabilities":["H.264Streaming","JPEGSnapshot","SignalStrength","Privacy","Standby",{"Resolutions":[{"text":"1080p","x":1920,"y":1088},{"text":"720p","x":1280,"y":720},{"text":"480p","x":848,"y":480},{"text":"360p","x":640,"y":352},{"text":"240p","x":416,"y":240}]},{"TimedStreamDuration":{"min":5,"max":120,"default":10}},{"TriggerEndStreamDuration":{"min":5,"max":300,"default":300}},{"Actions":[{"recordVideo":[{"StopActions":["timeout","triggerEndDetected"]}]},"sendEmailAlert","pushNotification"]},{"Triggers":[{"type":"pirMotionActive","sensitivity":{"type":"integer","min":1,"max":100,"step":1,"default":80}}]}],"zoom":{"topleftx":0,"toplefty":0,"bottomrightx":1280,"bottomrighty":720},"mic":{"mute":false,"volume":100},"speaker":{"mute":false,"volume":100},"streamingMode":"eventBased","continuousStreamState":"inactive","motion":{"sensitivity":5,"zones":[]},"resolution":{"width":1280,"height":720},"idleLedEnable":true,"privacyActive":false,"standbyActive":false,"connectionState":"available","activityState":"idle","swVersion":"1.2.16720","hwVersion":"H7","modelId":"VMC3030","motionSetupModeEnabled":false,"motionSetupModeSensitivity":80,"motionDetected":false,"audioDetected":false,"hasStreamed":true,"olsonTimeZone":"America/Chicago","name":"","nightVisionMode":1},{"interfaceVersion":3,"serialNumber":"48B4597FD9B8E","batteryLevel":0,"signalStrength":4,"brightness":0,"mirror":false,"flip":false,"powerSaveMode":3,"capabilities":["H.264Streaming","JPEGSnapshot","SignalStrength","Privacy","Standby",{"Resolutions":[{"text":"1080p","x":1920,"y":1088},{"text":"720p","x":1280,"y":720},{"text":"480p","x":848,"y":480},{"text":"360p","x":640,"y":352},{"text":"240p","x":416,"y":240}]},{"TimedStreamDuration":{"min":5,"max":120,"default":10}},{"TriggerEndStreamDuration":{"min":5,"max":300,"default":300}},{"Actions":[{"recordVideo":[{"StopActions":["timeout","triggerEndDetected"]}]},"sendEmailAlert","pushNotification"]},{"Triggers":[{"type":"pirMotionActive","sensitivity":{"type":"integer","min":1,"max":100,"step":1,"default":80}}]}],"zoom":{"topleftx":0,"toplefty":0,"bottomrightx":1280,"bottomrighty":720},"mic":{"mute":false,"volume":100},"speaker":{"mute":false,"volume":100},"streamingMode":"eventBased","continuousStreamState":"inactive","motion":{"sensitivity":5,"zones":[]},"resolution":{"width":1280,"height":720},"idleLedEnable":true,"privacyActive":false,"standbyActive":false,"connectionState":"batteryCritical","activityState":"idle","swVersion":"1.2.16720","hwVersion":"H7","modelId":"VMC3030","motionSetupModeEnabled":false,"motionSetupModeSensitivity":80,"motionDetected":false,"audioDetected":false,"hasStreamed":true,"olsonTimeZone":"America/Chicago","name":"","nightVisionMode":1},{"interfaceVersion":3,"serialNumber":"48B4597VD8FF5","batteryLevel":0,"signalStrength":4,"brightness":2,"mirror":true,"flip":true,"powerSaveMode":3,"capabilities":["H.264Streaming","JPEGSnapshot","SignalStrength","Privacy","Standby",{"Resolutions":[{"text":"1080p","x":1920,"y":1088},{"text":"720p","x":1280,"y":720},{"text":"480p","x":848,"y":480},{"text":"360p","x":640,"y":352},{"text":"240p","x":416,"y":240}]},{"TimedStreamDuration":{"min":5,"max":120,"default":10}},{"TriggerEndStreamDuration":{"min":5,"max":300,"default":300}},{"Actions":[{"recordVideo":[{"StopActions":["timeout","triggerEndDetected"]}]},"sendEmailAlert","pushNotification"]},{"Triggers":[{"type":"pirMotionActive","sensitivity":{"type":"integer","min":1,"max":100,"step":1,"default":80}}]}],"zoom":{"topleftx":0,"toplefty":0,"bottomrightx":1280,"bottomrighty":720},"mic":{"mute":false,"volume":100},"speaker":{"mute":false,"volume":100},"streamingMode":"eventBased","continuousStreamState":"inactive","motion":{"sensitivity":5,"zones":[]},"resolution":{"width":1280,"height":720},"idleLedEnable":true,"privacyActive":false,"standbyActive":false,"connectionState":"batteryCritical","activityState":"idle","swVersion":"1.2.16720","hwVersion":"H7","modelId":"VMC3030","motionSetupModeEnabled":false,"motionSetupModeSensitivity":80,"motionDetected":false,"audioDetected":false,"hasStreamed":true,"olsonTimeZone":"America/Chicago","name":"","nightVisionMode":1}],"action":"is","transId":"web!2dc849b8.9ffc2!1520148144127","to":"336-4764296_web","from":"48935B7SA9847"}

{"resource":"modes","properties":{"active":"mode1","modes":[{"name":"","type":"disarmed","rules":[],"id":"mode0"},{"name":"","type":"armed","rules":["rule5","rule3","rule0"],"id":"mode1"},{"name":"*****_DEFAULT_MODE_ARMED_*****","rules":["rule1"],"id":"mode2"},{"name":"Test Mode","rules":["rule6"],"id":"mode3"}]},"action":"is","transId":"web!bbb0ff1f.3c85f!1520148144127","to":"336-4764296_web","from":"48935B7SA9847"}


{"resource":"rules","properties":{"rules":[{"name":"Push notification if Front Door detects motion","protected":true,"triggers":[{"deviceId":"48B45974D8E54","sensitivity":80,"type":"pirMotionActive"}],"actions":[{"type":"recordVideo","deviceId":"48B45974D8E54","stopCondition":{"type":"timeout","timeout":120}},{"type":"pushNotification"}],"id":"rule0"},{"name":"Record camera (Back Patio) on motion.","protected":false,"triggers":[{"type":"pirMotionActive","deviceId":"48B4597VD8FF5","sensitivity":80}],"actions":[{"deviceId":"48B4597VD8FF5","type":"recordVideo","stopCondition":{"type":"timeout","timeout":10}},{"type":"pushNotification"}],"id":"rule1"},{"name":"Push notification if Inside detects motion","protected":true,"triggers":[{"deviceId":"48B4597FD9B8E","sensitivity":90,"type":"pirMotionActive"}],"actions":[{"deviceId":"48B4597FD9B8E","type":"recordVideo","stopCondition":{"type":"timeout","timeout":120}}],"id":"rule3"},{"name":"Push notification if Back Patio detects motion","protected":true,"triggers":[{"deviceId":"48B4597VD8FF5","sensitivity":100,"type":"pirMotionActive"}],"actions":[{"deviceId":"48B4597VD8FF5","type":"recordVideo","stopCondition":{"type":"timeout","timeout":30}},{"type":"pushNotification"}],"id":"rule5"},{"name":"Push notification & Email alert if Back Patio detects motion","protected":false,"triggers":[{"type":"pirMotionActive","deviceId":"48B4597VD8FF5","sensitivity":80}],"actions":[{"type":"sendEmailAlert","recipients":["__OWNER_EMAIL__"]},{"type":"pushNotification"}],"id":"rule6"}]},"action":"is","transId":"web!bff59099.cbd6d!1520148144127","to":"336-4764296_web","from":"48935B7SA9847"}


{"resource":"subscriptions/336-4764296_web","transId":"web!ddda6350.ba92c!1520148172685","action":"is","from":"48935B7SA9847","to":"336-4764296_web","properties":{"devices":["48935B7SA9847"],"url":"https://vzweb05-prod.vz.netgear.com/hmsweb/publish/48935B7SA9847/336-4764296/37da66eb-023f-4965-bb8b-480687881b65"}}


{"resource":"subscriptions/336-4764296_web","transId":"web!d5739e5.077af!1520148202738","action":"is","from":"48935B7SA9847","to":"336-4764296_web","properties":{"devices":["48935B7SA9847"],"url":"https://vzweb05-prod.vz.netgear.com/hmsweb/publish/48935B7SA9847/336-4764296/7d9cc5d7-a908-4f22-aaaa-dbb70c8616d6"}}
*/
