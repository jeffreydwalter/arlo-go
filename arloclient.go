package arloclient

import (
	"time"

	"github.com/jeffreydwalter/arloclient/internal/request"
	"github.com/jeffreydwalter/arloclient/internal/util"

	"github.com/pkg/errors"
)

type Arlo struct {
	user    string
	pass    string
	client  *request.Client
	Account *Account
	Devices *Devices
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

func Login(user string, pass string) (*Arlo, error) {

	a := newArlo(user, pass)

	body := map[string]string{"email": a.user, "password": a.pass}
	resp, err := a.client.Post(LoginUri, body, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "login request failed")
	}

	var loginResponse LoginResponse
	if err := util.Decode(resp.ParsedBody, &loginResponse); err != nil {
		return nil, err
	}

	if loginResponse.Success {
		// Cache the auth token.
		a.client.BaseHttpHeader.Add("Authorization", loginResponse.Data.Token)

		// Save the account info with the Arlo struct.
		a.Account = &loginResponse.Data

		if deviceResponse, err := a.GetDevices(); err != nil {
			return nil, err
		} else {
			if !deviceResponse.Success {
				return nil, err
			}
			a.Devices = &deviceResponse.Data
		}
	} else {
		return nil, errors.New("failed to login")
	}

	return a, nil
}

func (a *Arlo) Logout() (*Status, error) {

	resp, err := a.client.Put(LogoutUri, nil, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "logout request failed")
	}

	var status Status
	if err := util.Decode(resp.ParsedBody, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (a *Arlo) GetDevices() (*DeviceResponse, error) {

	resp, err := a.client.Get(DevicesUri, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "get devices request failed")
	}

	var deviceResponse DeviceResponse
	if err := util.Decode(resp.ParsedBody, &deviceResponse); err != nil {
		return nil, err
	}

	return &deviceResponse, nil
}

func (a *Arlo) GetLibraryMetaData(fromDate, toDate time.Time) (*LibraryMetaDataResponse, error) {

	body := map[string]string{"dateFrom": fromDate.Format("20060102"), "dateTo": toDate.Format("20060102")}
	resp, err := a.client.Post(LibraryMetadataUri, body, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to get library metadata")
	}

	var libraryMetaDataResponse LibraryMetaDataResponse
	if err := util.Decode(resp.ParsedBody, &libraryMetaDataResponse); err != nil {
		return nil, err
	}

	return &libraryMetaDataResponse, nil
}

func (a *Arlo) GetLibrary(fromDate, toDate time.Time) (*LibraryResponse, error) {

	body := map[string]string{"dateFrom": fromDate.Format("20060102"), "dateTo": toDate.Format("20060102")}
	resp, err := a.client.Post(LibraryUri, body, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get library")
	}

	var libraryResponse LibraryResponse
	if err := util.Decode(resp.ParsedBody, &libraryResponse); err != nil {
		return nil, err
	}

	return &libraryResponse, nil
}

func (a *Arlo) UpdateDeviceName(d Device, name string) (*Status, error) {

	body := map[string]string{"deviceId": d.DeviceId, "deviceName": name, "parentId": d.ParentId}
	resp, err := a.client.Put(DeviceRenameUri, body, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to update device name")
	}

	var status Status
	if err := util.Decode(resp.ParsedBody, &status); err != nil {
		return nil, err
	}

	return &status, nil

	return nil, errors.New("Device not found")
}

// UpdateProfile takes a first and last name, and updates the user profile with that information.
func (a *Arlo) UpdateProfile(firstName, lastName string) (*Status, error) {

	body := map[string]string{"firstName": firstName, "lastName": lastName}
	resp, err := a.client.Put(UserProfileUri, body, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to update profile")
	}

	var status Status
	if err := util.Decode(resp.ParsedBody, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (a *Arlo) UpdatePassword(password string) (*Status, error) {

	body := map[string]string{"currentPassword": a.pass, "newPassword": password}
	resp, err := a.client.Post(UserChangePasswordUri, body, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update password")
	}

	var status Status
	if err := util.Decode(resp.ParsedBody, &status); err != nil {
		return nil, err
	}

	if status.Success {
		a.pass = password
	}

	return &status, nil
}

/*
 This is an example of the json you would pass in the body to UpdateFriends():
{
  "firstName":"Some",
  "lastName":"Body",
  "devices":{
    "XXXXXXXXXXXXX":"Camera 1",
    "XXXXXXXXXXXXX":"Camera 2 ",
    "XXXXXXXXXXXXX":"Camera 3"
  },
  "lastModified":1463977440911,
  "adminUser":true,
  "email":"user@example.com",
  "id":"XXX-XXXXXXX"
}
*/
func (a *Arlo) UpdateFriends(f Friend) (*Status, error) {

	resp, err := a.client.Put(UserFriendsUri, f, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update friends")
	}

	var status Status
	if err := util.Decode(resp.ParsedBody, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (a *Arlo) UpdateDisplayOrder(d DeviceOrder) (*Status, error) {

	resp, err := a.client.Post(DeviceDisplayOrderUri, d, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update display order")
	}

	var status Status
	if err := util.Decode(resp.ParsedBody, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

/*
##
# Delete a single video recording from Arlo.
#
# All of the date info and device id you need to pass into this method are given in the results of the GetLibrary() call.
#
##
*/
func (a *Arlo) DeleteRecording(r *Recording) (*Status, error) {

	body := map[string]map[string]interface{}{"data": {"createdDate": r.CreatedDate, "utcCreatedDate": r.UtcCreatedDate, "deviceId": r.DeviceId}}
	resp, err := a.client.Post(LibraryRecycleUri, body, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to delete recording")
	}

	var status Status
	if err := util.Decode(resp.ParsedBody, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

/*
##
# Delete a batch of video recordings from Arlo.
#
# The GetLibrary() call response json can be passed directly to this method if you'd like to delete the same list of videos you queried for.
# If you want to delete some other batch of videos, then you need to send an array of objects representing each video you want to delete.
#
#[
#  {
#    "createdDate":"20160904",
#    "utcCreatedDate":1473010280395,
#    "deviceId":"XXXXXXXXXXXXX"
#  },
#  {
#    "createdDate":"20160904",
#    "utcCreatedDate":1473010280395,
#    "deviceId":"XXXXXXXXXXXXX"
#  }
#]
##
func (a *Arlo) BatchDeleteRecordings(recording_metadata):
return a.client.Post('https://arlo.netgear.com/hmsweb/users/library/recycle', {'data':recording_metadata})

##
# Returns the whole video from the presignedContentUrl.
#
# Obviously, this function is generic and could be used to download anything. :)
##
func (a *Arlo) GetRecording(url, chunk_size=4096):
video = ''
r = requests.get(url, stream=True)
r.raise_for_status()

for chunk in r.iter_content(chunk_size):
if chunk: video += chunk
return video


##
# This function returns a json object containing the rtmps url to the requested video stream.
# You will need the to install a library to handle streaming of this protocol: https://pypi.python.org/pypi/python-librtmp
#
# The request to /users/devices/startStream returns:
#{ "url":"rtmps://vzwow09-z2-prod.vz.netgear.com:80/vzmodulelive?egressToken=b1b4b675_ac03_4182_9844_043e02a44f71&userAgent=web&cameraId=48B4597VD8FF5_1473010750131" }
#
##
func (a *Arlo) StartStream(camera):
return a.client.Post('https://arlo.netgear.com/hmsweb/users/devices/startStream', {"to":camera.get('parentId'),"from":self.user_id+"_web","resource":"cameras/"+camera.get('deviceId'),"action":"set","publishResponse":True,"transId":self.genTransId(),"properties":{"activityState":"startUserStream","cameraId":camera.get('deviceId')}}, headers={"xcloudId":camera.get('xCloudId')})

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
