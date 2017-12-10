package arlo

import (
	"github.com/pkg/errors"
)

// Account is the account data.
type Account struct {
	UserId        string  `json:"userId"`
	Email         string  `json:"email"`
	Token         string  `json:"token"`
	PaymentId     string  `json:"paymentId"`
	Authenticated uint32  `json:"authenticated"`
	AccountStatus string  `json:"accountStatus"`
	SerialNumber  string  `json:"serialNumber"`
	CountryCode   string  `json:"countryCode"`
	TocUpdate     bool    `json:"tocUpdate"`
	PolicyUpdate  bool    `json:"policyUpdate"`
	ValidEmail    bool    `json:"validEmail"`
	Arlo          bool    `json:"arlo"`
	DateCreated   float64 `json:"dateCreated"`
}

type Friend struct {
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Devices      DeviceOrder `json:"devices"`
	LastModified float64     `json:"lastModified"`
	AdminUser    bool        `json:"adminUser"`
	Email        string      `json:"email"`
	Id           string      `json:"id"`
}

func Login(user string, pass string) (*Arlo, error) {

	a := newArlo(user, pass)

	body := map[string]string{"email": a.user, "password": a.pass}
	resp, err := a.client.Post(LoginUri, body, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "login request failed")
	}

	var loginResponse LoginResponse
	if err := resp.Decode(&loginResponse); err != nil {
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

			// Cache the devices as their respective types.
			a.Basestations = deviceResponse.Data.Basestations()
			a.Cameras = deviceResponse.Data.Cameras()

			// Set the XCloudId header for future requests. You can override this on a per-request basis if needed.
			a.client.BaseHttpHeader.Add("xCloudId", deviceResponse.Data[0].XCloudId)
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
	if err := resp.Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}

// UpdateProfile takes a first and last name, and updates the user profile with that information.
func (a *Arlo) UpdateProfile(firstName, lastName string) (*Status, error) {

	body := map[string]string{"firstName": firstName, "lastName": lastName}
	resp, err := a.client.Put(UserProfileUri, body, nil)

	if err != nil {
		return nil, errors.WithMessage(err, "failed to update profile")
	}

	var status Status
	if err := resp.Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (a *Arlo) UpdatePassword(pass string) (*Status, error) {

	body := map[string]string{"currentPassword": a.pass, "newPassword": pass}
	resp, err := a.client.Post(UserChangePasswordUri, body, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to update password")
	}

	var status Status
	if err := resp.Decode(&status); err != nil {
		return nil, err
	}

	if status.Success {
		a.pass = pass
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
	if err := resp.Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}
