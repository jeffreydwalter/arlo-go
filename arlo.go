package arlo

import (
	"github.com/jeffreydwalter/arlo/internal/request"
)

type Arlo struct {
	user         string
	pass         string
	client       *request.Client
	Account      *Account
	Basestations *Basestations
	Cameras      *Cameras
}

func newArlo(user string, pass string) *Arlo {

	c, _ := request.NewClient(BaseUrl)
	arlo := &Arlo{
		user:   user,
		pass:   pass,
		client: c,
	}

	return arlo
}

/*
##
# This function causes the camera to record a snapshot.
#
# You can get the timezone from GetDevices().
##
func (a *Arlo) TakeSnapshot(camera):
stream_url = self.StartStream(camera)
a.client.Post('https://arlo.netgear.com/hmsweb/users/devices/takeSnapshot', {'xcloudId':camera.get('xCloudId'),'parentId':camera.get('parentId'),'deviceId':camera.get('deviceId'),'olsonTimeZone':camera.get('properties', {}).get('olsonTimeZone')}, headers={"xcloudId":camera.get('xCloudId')})
return stream_url;

##
# This function causes the camera to start recording.
#
# You can get the timezone from GetDevices().
##
func (a *Arlo) StartRecording(camera):
stream_url = self.StartStream(camera)
a.client.Post('https://arlo.netgear.com/hmsweb/users/devices/startRecord', {'xcloudId':camera.get('xCloudId'),'parentId':camera.get('parentId'),'deviceId':camera.get('deviceId'),'olsonTimeZone':camera.get('properties', {}).get('olsonTimeZone')}, headers={"xcloudId":camera.get('xCloudId')})
return stream_url

##
# This function causes the camera to stop recording.
#
# You can get the timezone from GetDevices().
##
func (a *Arlo) StopRecording(camera):
return a.client.Post('https://arlo.netgear.com/hmsweb/users/devices/stopRecord', {'xcloudId':camera.get('xCloudId'),'parentId':camera.get('parentId'),'deviceId':camera.get('deviceId'),'olsonTimeZone':camera.get('properties', {}).get('olsonTimeZone')}, headers={"xcloudId":camera.get('xCloudId')})
*/
