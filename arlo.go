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

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jeffreydwalter/arlo-go/internal/request"

	"github.com/pkg/errors"
)

type Arlo struct {
	user         string
	pass         string
	client       *request.Client
	Account      Account
	Basestations Basestations
	Cameras      Cameras
	rwmutex      sync.RWMutex
}

func newArlo(user string, pass string) (arlo *Arlo) {

	// Add important headers.
	baseHeaders := make(http.Header)
	baseHeaders.Add("DNT", "1")
	baseHeaders.Add("schemaVersion", "1")
	baseHeaders.Add("Host", "my.arlo.com")
	baseHeaders.Add("Referer", "https://my.arlo.com/")

	c, _ := request.NewClient(BaseUrl, baseHeaders)

	return &Arlo{
		user:   user,
		pass:   pass,
		client: c,
	}
}

func Login(user string, pass string) (arlo *Arlo, err error) {
	arlo = newArlo(user, pass)

	body := map[string]string{"email": arlo.user, "password": arlo.pass}
	resp, err := arlo.post(LoginV2Uri, "", body, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to login")
	}
	defer resp.Body.Close()

	var loginResponse LoginResponse
	if err := resp.Decode(&loginResponse); err != nil {
		return nil, err
	}

	if loginResponse.Success {
		// Cache the auth token.
		arlo.client.AddHeader("Authorization", loginResponse.Data.Token)

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

func (a *Arlo) CheckSession() (session *Session, err error) {
	msg := "failed to get session"
	resp, err := a.get(SessionUri, "", nil)
	if err != nil {
		return nil, errors.WithMessage(err, msg)
	}
	defer resp.Body.Close()

	var response SessionResponse
	if err := resp.Decode(&response); err != nil {
		return nil, err
	}

	if response.Success == false {
		return nil, errors.WithMessage(errors.New(response.Reason), msg)
	}
	return &response.Data, nil
}

// GetDevices returns an array of all devices.
// When you call Login, this method is called and all devices are cached in the arlo object.
func (a *Arlo) GetDevices() (devices *Devices, err error) {
	resp, err := a.get(fmt.Sprintf(DevicesUri, time.Now().Format("20060102")), "", nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get devices")
	}
	defer resp.Body.Close()

	var response DeviceResponse
	if err := resp.Decode(&response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New("failed to get devices")
	}

	if len(response.Data) == 0 {
		return nil, errors.New("no devices found")
	}

	// Cache a pointer to the arlo object with each device.
	for i := range response.Data {
		response.Data[i].arlo = a
	}

	// Disconnect all of the basestations from the EventStream.
	for i := range a.Basestations {
		if err := a.Basestations[i].Disconnect(); err != nil {
			return nil, errors.WithMessage(err, "failed to get devices")
		}
	}

	a.rwmutex.Lock()
	// Cache the devices as their respective types.
	a.Cameras = *response.Data.GetCameras()
	a.Basestations = *response.Data.GetBasestations()
	a.rwmutex.Unlock()

	// subscribe each basestation to the EventStream.
	for i := range a.Basestations {
		if err := a.Basestations[i].Subscribe(); err != nil {
			return nil, errors.WithMessage(err, "failed to get devices")
		}
	}

	return &response.Data, nil
}

// GetProfile returns the user profile for the currently logged in user.
func (a *Arlo) GetProfile() (profile *UserProfile, err error) {
	resp, err := a.get(ProfileUri, "", nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user profile")
	}
	defer resp.Body.Close()

	var response UserProfileResponse
	if err := resp.Decode(&response); err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New("failed to get user profile")
	}

	return &response.Data, nil
}

// UpdateDisplayOrder sets the display order according to the order defined in the DeviceOrder given.
func (a *Arlo) UpdateDisplayOrder(d DeviceOrder) error {
	resp, err := a.post(CameraOrderUri, "", d, nil)
	return checkRequest(resp, err, "failed to display order")
}

// UpdateProfile takes a first and last name, and updates the user profile with that information.
func (a *Arlo) UpdateProfile(firstName, lastName string) error {
	body := map[string]string{"firstName": firstName, "lastName": lastName}
	resp, err := a.put(ProfileUri, "", body, nil)
	return checkRequest(resp, err, "failed to update profile")
}

func (a *Arlo) UpdatePassword(pass string) error {
	body := map[string]string{"currentPassword": a.pass, "newPassword": pass}
	resp, err := a.post(UpdatePasswordUri, "", body, nil)
	if err := checkRequest(resp, err, "failed to update password"); err != nil {
		return err
	}

	a.pass = pass

	return nil
}

func (a *Arlo) UpdateFriends(f Friend) error {
	resp, err := a.put(FriendsUri, "", f, nil)
	return checkRequest(resp, err, "failed to update friends")
}
