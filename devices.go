package arloclient

// Device is the device data.
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

// Devices is an array of Device objects.
type Devices []Device

// DeviceOrder is a hash of # XXXXXXXXXXXXX is the device id of each camera. You can get this from GetDevices().
/*
{
  "devices":{
    "XXXXXXXXXXXXX":1,
    "XXXXXXXXXXXXX":2,
    "XXXXXXXXXXXXX":3
}
*/
type DeviceOrder struct {
	Devices map[string]int
}

func (ds *Devices) Find(deviceId string) *Device {
	for _, d := range *ds {
		if d.DeviceId == deviceId {
			return &d
		}
	}

	return nil
}

func (ds *Devices) BaseStations() *Devices {
	var basestations Devices
	for _, d := range *ds {
		if d.DeviceType == "basestation" {
			basestations = append(basestations, d)
		}
	}
	return &basestations
}

func (ds *Devices) Cameras() *Devices {
	var cameras Devices
	for _, d := range *ds {
		if d.DeviceType != "basestation" {
			cameras = append(cameras, d)
		}
	}
	return &cameras
}
