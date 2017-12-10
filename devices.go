package arlo

import (
	"github.com/pkg/errors"
)

// A Device is the device data, this can be a camera, basestation, arloq, etc.
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

// A Basestation is a Device that's not type "camera" (basestation, arloq, arloqs, etc.).
// This type is here just for semantics. Some methods explicitly require a device of a certain type.
type Basestation Device

// Basestations is an array of Basestation objects.
type Basestations []Basestation

// A Camera is a Device of type "camera".
// This type is here just for semantics. Some methods explicitly require a device of a certain type.
type Camera Device

// Cameras is an array of Camera objects.
type Cameras []Camera

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

// Basestations returns a Basestations object containing all devices that are NOT type "camera".
// I did this because some device types, like arloq, don't have a basestation.
// So, when interacting with them you must treat them like a basestation and a camera.
// Cameras also includes decvices of this type, so you can get the same data there or cast.
func (ds *Devices) Basestations() *Basestations {
	var basestations Basestations
	for _, d := range *ds {
		if d.DeviceType != "camera" {
			basestations = append(basestations, Basestation(d))
		}
	}
	return &basestations
}

// Cameras returns a Cameras object containing all devices that are of type "camera".
// I did this because some device types, like arloq, don't have a basestation.
// So, when interacting with them you must treat them like a basestation and a camera.
// Basestations also includes decvices of this type, so you can get the same data there or cast.
func (ds *Devices) Cameras() *Cameras {
	var cameras Cameras
	for _, d := range *ds {
		if d.DeviceType != "basestation" {
			cameras = append(cameras, Camera(d))
		}
	}
	return &cameras
}

// GetDevices returns an array of all devices.
// When you call Login, this method is called and all devices are cached in the Arlo object.
func (a *Arlo) GetDevices() (*DeviceResponse, error) {

	resp, err := a.client.Get(DevicesUri, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to get devices")
	}

	var deviceResponse DeviceResponse
	if err := resp.Decode(&deviceResponse); err != nil {
		return nil, err
	}

	if len(deviceResponse.Data) == 0 {
		return nil, errors.New("no devices found")
	}

	return &deviceResponse, nil
}

// UpdateDeviceName sets the name of the given device to the name argument.
func (a *Arlo) UpdateDeviceName(d Device, name string) (*Status, error) {

	body := map[string]string{"deviceId": d.DeviceId, "deviceName": name, "parentId": d.ParentId}
	resp, err := a.client.Put(DeviceRenameUri, body, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to update device name")
	}

	var status Status
	if err := resp.Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil

	return nil, errors.New("device not found")
}

// UpdateDisplayOrder sets the display order according to the order defined in the DeviceOrder given.
func (a *Arlo) UpdateDisplayOrder(d DeviceOrder) (*Status, error) {

	resp, err := a.client.Post(DeviceDisplayOrderUri, d, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update display order")
	}

	var status Status
	if err := resp.Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}

// StartStream returns a json object containing the rtmps url to the requested video stream.
// You will need the to install a library to handle streaming of this protocol: https://pypi.python.org/pypi/python-librtmp
//
// The request to /users/devices/startStream returns:
// NOTE: { "url":"rtmps://vzwow09-z2-prod.vz.netgear.com:80/vzmodulelive?egressToken=b1b4b675_ac03_4182_9844_043e02a44f71&userAgent=web&cameraId=48B4597VD8FF5_1473010750131" }
func (a *Arlo) StartStream(c Camera) (*StartStreamResponse, error) {

	var n Notification
	n.To = c.ParentId
	n.From = c.UserId
	n.Resource = "cameras/" + c.DeviceId
	n.Action = "set"
	n.PublishResponse = true
	n.TransId = ""
	n.Properties.ActivityState = "startUserStream"
	n.Properties.CameraId = c.DeviceId

	// {"to":camera.get('parentId'),"from":self.user_id+"_web","resource":"cameras/"+camera.get('deviceId'),"action":"set","publishResponse":True,"transId":self.genTransId(),"properties":{"activityState":"startUserStream","cameraId":camera.get('deviceId')}}, headers={"xcloudId":camera.get('xCloudId')}
	resp, err := a.client.Post(DeviceStartStreamUri, n, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to start stream")
	}

	var startStreamResponse StartStreamResponse
	if err := resp.Decode(&startStreamResponse); err != nil {
		return nil, err
	}

	return &startStreamResponse, nil
}
