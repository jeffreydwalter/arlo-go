package arlo

import (
	"github.com/jeffreydwalter/arlo-golang/internal/request"

	"github.com/pkg/errors"
)

type Arlo struct {
	user         string
	pass         string
	client       *request.Client
	Account      Account
	Basestations Basestations
	Cameras      Cameras
}

func newArlo(user string, pass string) (arlo *Arlo) {

	c, _ := request.NewClient(BaseUrl)

	// Add important headers.
	c.BaseHttpHeader.Add("DNT", "1")
	c.BaseHttpHeader.Add("schemaVersion", "1")
	c.BaseHttpHeader.Add("Host", "arlo.netgear.com")
	c.BaseHttpHeader.Add("Referer", "https://arlo.netgear.com/")

	return &Arlo{
		user:   user,
		pass:   pass,
		client: c,
	}
}

func Login(user string, pass string) (arlo *Arlo, err error) {
	arlo = newArlo(user, pass)

	body := map[string]string{"email": arlo.user, "password": arlo.pass}
	resp, err := arlo.post(LoginUri, "", body, nil)
	if err := checkHttpRequest(resp, err, "login request failed"); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var loginResponse LoginResponse
	if err := resp.Decode(&loginResponse); err != nil {
		return nil, err
	}

	if loginResponse.Success {
		// Cache the auth token.
		arlo.client.BaseHttpHeader.Add("Authorization", loginResponse.Data.Token)

		// Save the account info with the arlo struct.
		arlo.Account = loginResponse.Data

		// Get the devices, which also caches them on the arlo object.
		if _, err := arlo.GetDevices(); err != nil {
			return nil, errors.WithMessage(err, "failed to login")
		}
	} else {
		return nil, errors.New("failed to login")
	}

	return arlo, nil
}

func (a *Arlo) Logout() error {
	resp, err := a.put(LogoutUri, "", nil, nil)
	return checkRequest(resp, err, "failed to logout")
}

// GetDevices returns an array of all devices.
// When you call Login, this method is called and all devices are cached in the arlo object.
func (a *Arlo) GetDevices() (devices Devices, err error) {
	resp, err := a.get(DevicesUri, "", nil)
	if err := checkHttpRequest(resp, err, "failed to get devices"); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var deviceResponse DeviceResponse
	if err := resp.Decode(&deviceResponse); err != nil {
		return nil, err
	}

	if !deviceResponse.Success {
		return nil, errors.New("failed to get devices")
	}

	if len(deviceResponse.Data) == 0 {
		return nil, errors.New("no devices found")
	}

	for i := range deviceResponse.Data {
		deviceResponse.Data[i].arlo = a
	}

	// Unsubscribe all of the basestations from the EventStream.
	for i := range a.Basestations {
		if err := a.Basestations[i].Unsubscribe(); err != nil {
			return nil, errors.WithMessage(err, "failed to get devices")
		}
	}

	// Cache the devices as their respective types.
	a.Cameras = deviceResponse.Data.GetCameras()
	a.Basestations = deviceResponse.Data.GetBasestations()

	// Subscribe each basestation to the EventStream.
	for i := range a.Basestations {
		if err := a.Basestations[i].Subscribe(); err != nil {
			return nil, errors.WithMessage(err, "failed to get devices")
		}
	}

	return deviceResponse.Data, nil
}

// UpdateDisplayOrder sets the display order according to the order defined in the DeviceOrder given.
func (a *Arlo) UpdateDisplayOrder(d DeviceOrder) error {
	resp, err := a.post(DeviceDisplayOrderUri, "", d, nil)
	return checkRequest(resp, err, "failed to display order")
}

// UpdateProfile takes a first and last name, and updates the user profile with that information.
func (a *Arlo) UpdateProfile(firstName, lastName string) error {
	body := map[string]string{"firstName": firstName, "lastName": lastName}
	resp, err := a.put(UserProfileUri, "", body, nil)
	return checkRequest(resp, err, "failed to update profile")
}

func (a *Arlo) UpdatePassword(pass string) error {
	body := map[string]string{"currentPassword": a.pass, "newPassword": pass}
	resp, err := a.post(UserChangePasswordUri, "", body, nil)
	if err := checkRequest(resp, err, "failed to update password"); err != nil {
		return err
	}

	a.pass = pass

	return nil
}

func (a *Arlo) UpdateFriends(f Friend) error {
	resp, err := a.put(UserFriendsUri, "", f, nil)
	return checkRequest(resp, err, "failed to update friends")
}
