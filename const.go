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

const (
	BaseUrl               = "https://arlo.netgear.com/hmsweb"
	LoginUri              = "/login/v2"
	LogoutUri             = "/logout"
	SubscribeUri          = "/client/subscribe?token=%s"
	UnsubscribeUri        = "/client/unsubscribe"
	NotifyUri             = "/users/devices/notify/%s"
	ServiceLevelUri       = "/users/serviceLevel"
	OffersUri             = "/users/payment/offers"
	UserProfileUri        = "/users/profile"
	PushToTalkUri         = "/users/devices/%s/pushtotalk"
	UserChangePasswordUri = "/users/changePassword"
	UserSessionUri        = "/users/session"
	UserFriendsUri        = "/users/friends"
	UserLocationsUri      = "/users/locations"
	UserLocationUri       = "/users/locations/%s"
	LibraryUri            = "/users/library"
	LibraryMetadataUri    = "/users/library/metadata"
	LibraryRecycleUri     = "/users/library/recycle"
	LibraryResetUri       = "/users/library/reset"
	DevicesUri            = "/users/devices"
	DeviceRenameUri       = "/users/devices/renameDevice"
	DeviceDisplayOrderUri = "/users/devices/displayOrder"
	DeviceTakeSnapshotUri = "/users/devices/takeSnapshot"
	DeviceStartRecordUri  = "/users/devices/startRecord"
	DeviceStopRecordUri   = "/users/devices/stopRecord"
	DeviceStartStreamUri  = "/users/devices/startStream"

	DeviceTypeBasestation = "basestation"
	DeviceTypeCamera      = "camera"
	DeviceTypeArloQ       = "arloq"
	TransIdPrefix         = "web"
)
