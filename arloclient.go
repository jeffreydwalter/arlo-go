package arloclient

import (
	"log"
	"time"

	"github.com/jeffreydwalter/arloclient/internal/request"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Arlo struct {
	user    string
	pass    string
	client  *request.Client
	account Account
}

func NewArlo(user string, pass string) (*Arlo, error) {

	c, _ := request.NewClient(BaseUrl)
	arlo := &Arlo{
		user:   user,
		pass:   pass,
		client: c,
	}

	if _, err := arlo.Login(); err != nil {
		return nil, errors.WithMessage(err, "failed to create arlo object")
	}

	return arlo, nil
}

func (a *Arlo) Login() (*Account, error) {

	resp, err := a.client.Post(LoginUri, Credentials{Email: a.user, Password: a.pass}, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to login")
	}

	var loginResponse LoginResponse
	if err := mapstructure.Decode(resp.Data, &loginResponse); err != nil {
		return nil, errors.Wrap(err, "failed to create loginresponse object")
	}

	if !loginResponse.Success {
		return nil, errors.New("request was unsuccessful")
	}

	// Cache the auth token.
	a.client.BaseHttpHeader.Add("Authorization", loginResponse.Data.Token)

	// Save the account info with the Arlo struct.
	a.account = loginResponse.Data

	return &loginResponse.Data, nil
}

func (a *Arlo) Logout() (*request.Response, error) {

	return a.client.Put(LogoutUri, nil, nil)
}

func (a *Arlo) GetDevices() (*Devices, error) {

	resp, err := a.client.Get(DevicesUri, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to get devices")
	}

	var deviceResponse DeviceResponse
	if err := mapstructure.Decode(resp.Data, &deviceResponse); err != nil {
		return nil, errors.Wrap(err, "failed to create deviceresponse object")
	}

	if !deviceResponse.Success {
		return nil, errors.New("request was unsuccessful")
	}

	return &deviceResponse.Data, nil
}

func (a *Arlo) GetLibraryMetaData(fromDate, toDate time.Time) (*LibraryMetaData, error) {

	resp, err := a.client.Post(LibraryMetadataUri, Duration{fromDate.Format("20060102"), toDate.Format("20060102")}, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to get library metadata")
	}

	log.Printf("GETLIBRARYMETADATA: %v", resp.Data)

	var libraryMetaDataResponse LibraryMetaDataResponse
	if err := mapstructure.Decode(resp.Data, &libraryMetaDataResponse); err != nil {
		return nil, errors.WithMessage(err, "failed to create librarymetadataresponse object")
	}

	if !libraryMetaDataResponse.Success {
		return nil, errors.New("request was unsuccessful")
	}

	return &libraryMetaDataResponse.Data, nil
}

func (a *Arlo) UpdateProfile(firstName, lastName string) (*UserProfile, error) {

	resp, err := a.client.Put(UserProfileUri, FullName{firstName, lastName}, nil)

	if err != nil {
		return nil, err
	}

	var userProfileResponse UserProfileResponse
	if err := mapstructure.Decode(resp.Data, &userProfileResponse); err != nil {
		return nil, err
	}

	if !userProfileResponse.Success {
		return nil, err
	}

	return &userProfileResponse.Data, nil
}

func (a *Arlo) UpdatePassword(password string) error {

	_, err := a.client.Post(UserChangePasswordUri, PasswordPair{a.pass, password}, nil)
	if err != nil {
		a.pass = password
	}
	return err
}

/*
##
# This is an example of the json you would pass in the body to UpdateFriends():
#{
#  "firstName":"Some",
#  "lastName":"Body",
#  "devices":{
#    "XXXXXXXXXXXXX":"Camera 1",
#    "XXXXXXXXXXXXX":"Camera 2 ",
#    "XXXXXXXXXXXXX":"Camera 3"
#  },
#  "lastModified":1463977440911,
#  "adminUser":true,
#  "email":"user@example.com",
#  "id":"XXX-XXXXXXX"
#}
##
func (a *Arlo) UpdateFriends(body):
return a.client.Put('https://arlo.netgear.com/hmsweb/users/friends', body)

func (a *Arlo) UpdateDeviceName(device, name):
return a.client.Put('https://arlo.netgear.com/hmsweb/users/devices/renameDevice', {'deviceId':device.get('deviceId'), 'deviceName':name, 'parentId':device.get('parentId')})

##
# This is an example of the json you would pass in the body to UpdateDisplayOrder() of your devices in the UI.
#
# XXXXXXXXXXXXX is the device id of each camera. You can get this from GetDevices().
#{
#  "devices":{
#    "XXXXXXXXXXXXX":1,
#    "XXXXXXXXXXXXX":2,
#    "XXXXXXXXXXXXX":3
#  }
#}
##
func (a *Arlo) UpdateDisplayOrder(body):
return a.client.Post('https://arlo.netgear.com/hmsweb/users/devices/displayOrder', body)

##
# This call returns the following:
# presignedContentUrl is a link to the actual video in Amazon AWS.
# presignedThumbnailUrl is a link to the thumbnail .jpg of the actual video in Amazon AWS.
#
#[
# {
#  "mediaDurationSecond": 30,
#  "contentType": "video/mp4",
#  "name": "XXXXXXXXXXXXX",
#  "presignedContentUrl": "https://arlos3-prod-z2.s3.amazonaws.com/XXXXXXX_XXXX_XXXX_XXXX_XXXXXXXXXXXXX/XXX-XXXXXXX/XXXXXXXXXXXXX/recordings/XXXXXXXXXXXXX.mp4?AWSAccessKeyId=XXXXXXXXXXXXXXXXXXXX&Expires=1472968703&Signature=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
#  "lastModified": 1472881430181,
#  "localCreatedDate": XXXXXXXXXXXXX,
#  "presignedThumbnailUrl": "https://arlos3-prod-z2.s3.amazonaws.com/XXXXXXX_XXXX_XXXX_XXXX_XXXXXXXXXXXXX/XXX-XXXXXXX/XXXXXXXXXXXXX/recordings/XXXXXXXXXXXXX_thumb.jpg?AWSAccessKeyId=XXXXXXXXXXXXXXXXXXXX&Expires=1472968703&Signature=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
#  "reason": "motionRecord",
#  "deviceId": "XXXXXXXXXXXXX",
#  "createdBy": "XXXXXXXXXXXXX",
#  "createdDate": "20160903",
#  "timeZone": "America/Chicago",
#  "ownerId": "XXX-XXXXXXX",
#  "utcCreatedDate": XXXXXXXXXXXXX,
#  "currentState": "new",
#  "mediaDuration": "00:00:30"
# }
#]
##
func (a *Arlo) GetLibrary(from_date, to_date):
return a.client.Post('https://arlo.netgear.com/hmsweb/users/library', {'dateFrom':from_date, 'dateTo':to_date})

##
# Delete a single video recording from Arlo.
#
# All of the date info and device id you need to pass into this method are given in the results of the GetLibrary() call.
#
##
func (a *Arlo) DeleteRecording(camera, created_date, utc_created_date):
return a.client.Post('https://arlo.netgear.com/hmsweb/users/library/recycle', {'data':[{'createdDate':created_date,'utcCreatedDate':utc_created_date,'deviceId':camera.get('deviceId')}]})

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
