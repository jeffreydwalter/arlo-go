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

// A Device is the device data, this can be a camera, basestation, arloq, etc.
type Device struct {
	arlo                          *Arlo        // Let's hold a reference to the parent arlo object since it holds the http.Client object and references to all devices.
	AnalyticsEnabled              bool         `json:"analyticsEnabled"`
	ArloMobilePlan                bool         `json:"arloMobilePlan"`
	ArloMobilePlanId              string       `json:"arloMobilePlanId"`
	ArloMobilePlanName            string       `json:"arloMobilePlanName"`
	ArloMobilePlanThreshold       int          `json:"arloMobilePlanThreshold"`
	Connectivity                  Connectivity `json:"connectivity"`
	CriticalBatteryState          bool         `json:"criticalBatteryState"`
	DateCreated                   int64        `json:"dateCreated"`
	DeviceId                      string       `json:"deviceId"`
	DeviceName                    string       `json:"deviceName"`
	DeviceType                    string       `json:"deviceType"`
	DisplayOrder                  uint8        `json:"displayOrder"`
	FirmwareVersion               string       `json:"firmwareVersion"`
	InterfaceVersion              string       `json:"interfaceVersion"`
	InterfaceSchemaVer            string       `json:"interfaceSchemaVer"`
	LastImageUploaded             string       `json:"lastImageUploaded"`
	LastModified                  int64        `json:"lastModified"`
	MigrateActivityZone           bool         `json:"migrateActivityZone"`
	MobileCarrier                 string       `json:"mobileCarrier"`
	MobileTrialUsed               bool         `json:"mobileTrialUsed"`
	PermissionsFilePath           string       `json:"permissionsFilePath"`
	PermissionsSchemaVer          string       `json:"permissionsSchemaVer"`
	PermissionsVerison            string       `json:"permissionsVerison"` // WTF? Netgear developers think this is OK... *sigh*
	PermissionsVersion            string       `json:"permissionsVersion"`
	PresignedFullFrameSnapshotUrl string       `json:"presignedFullFrameSnapshotUrl"`
	PresignedLastImageUrl         string       `json:"presignedLastImageUrl"`
	PresignedSnapshotUrl          string       `json:"presignedSnapshotUrl"`
	MediaObjectCount              uint8        `json:"mediaObjectCount"`
	ModelId                       string       `json:"modelId"`
	Owner                         Owner        `json:"owner"`
	ParentId                      string       `json:"parentId"`
	Properties                    Properties   `json:"properties"`
	UniqueId                      string       `json:"uniqueId"`
	UserId                        string       `json:"userId"`
	UserRole                      string       `json:"userRole"`
	State                         string       `json:"state"`
	XCloudId                      string       `json:"xCloudId"`
}

// Devices is a slice of Device objects.
type Devices []Device

// A DeviceOrder holds a map of device ids and a numeric index. The numeric index is the device order.
// Device order is mainly used by the UI to determine which order to show the devices.
/*
{
  "devices":{
    "XXXXXXXXXXXXX":1,
    "XXXXXXXXXXXXX":2,
    "XXXXXXXXXXXXX":3
}
*/
type DeviceOrder struct {
	Devices map[string]int `json:"devices"`
}

// Find returns a device with the device id passed in.
func (ds *Devices) Find(deviceId string) *Device {
	for _, d := range *ds {
		if d.DeviceId == deviceId {
			return &d
		}
	}

	return nil
}

func (ds Devices) FindCameras(basestationId string) Cameras {
	cs := new(Cameras)
	for _, d := range ds {
		if d.ParentId == basestationId {
			*cs = append(*cs, Camera(d))
		}
	}

	return *cs
}

func (d Device) IsBasestation() bool {
	return d.DeviceType == DeviceTypeBasestation || d.DeviceId == d.ParentId
}

func (d Device) IsCamera() bool {
	switch(d.DeviceType) {
        case
            DeviceTypeCamera,
            DeviceTypeArloQ:
        return true
    }
    return false
}

func (d Device) IsArloQ() bool {
	return d.DeviceType == DeviceTypeArloBridge
}

func (d Device) IsLight() bool {
	return d.DeviceType == DeviceTypeLights
}

func (d Device) IsSiren() bool {
	return d.DeviceType == DeviceTypeSiren
}

// GetBasestations returns a Basestations object containing all devices that are NOT type "camera".
// I did this because some device types, like arloq, don't have a basestation.
// So, when interacting with them you must treat them like a basestation and a camera.
// Cameras also includes devices of this type, so you can get the same data there or cast.
func (ds Devices) GetBasestations() *Basestations {
	basestations := new(Basestations)
	for _, d := range ds {
		if d.IsBasestation() || !d.IsCamera() {
			*basestations = append(*basestations, Basestation{Device: d})
		}
	}
	return basestations
}

// GetCameras returns a Cameras object containing all devices that are of type "camera".
// I did this because some device types, like arloq, don't have a basestation.
// So, when interacting with them you must treat them like a basestation and a camera.
// Basestations also includes devices of this type, so you can get the same data there or cast.
func (ds Devices) GetCameras() *Cameras {
	cameras := new(Cameras)
	for _, d := range ds {
		if d.IsCamera() || !d.IsBasestation() {
			*cameras = append(*cameras, Camera(d))
		}
	}
	return cameras
}

// UpdateDeviceName sets the name of the given device to the name argument.
func (d *Device) UpdateDeviceName(name string) error {
	body := map[string]string{"deviceId": d.DeviceId, "deviceName": name, "parentId": d.ParentId}
	resp, err := d.arlo.put(RenameDeviceUri, d.XCloudId, body, nil)
	return checkRequest(resp, err, "failed to update device name")
}
