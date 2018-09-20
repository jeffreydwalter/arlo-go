package arlo

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

/*
The methods in this file are all related to Arlo Baby (afaik).
They may apply to other camera types that have audio playback or nightlight capabilities.
*/

/*
 The follow methods are all related to the audio features of Arlo Baby.
*/

// SetVolume sets the volume of the audio playback to a level from 0-100.
func (c *Camera) SetVolume(volume int) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: SpeakerProperties{
			Speaker: VolumeProperties{
				Mute:   false,
				Volume: volume,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to set audio volume"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

// Mute mutes the audio playback.
func (c *Camera) Mute() (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: SpeakerProperties{
			Speaker: VolumeProperties{
				Mute: true,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to mute audio"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

// UnMute un-mutes the audio playback.
func (c *Camera) UnMute() (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: SpeakerProperties{
			Speaker: VolumeProperties{
				Mute: false,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to un-mute audio"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

// Play plays an audio track, specified by the track ID, from a given position starting from 0 seconds.
func (c *Camera) Play(trackId string, position int) error {

	// Defaulting to 'hugh little baby', which is a supplied track. Hopefully, the ID is the same for everyone.
	if trackId == "" {
		trackId = "2391d620-e491-4412-99f6-e9a40d6046ed"
	}

	if position < 0 {
		position = 0
	}

	payload := EventStreamPayload{
		Action:          "playTrack",
		Resource:        "audioPlayback/player",
		PublishResponse: false,
		Properties:      PlayTrackProperties{trackId, position},
		From:            fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:              c.ParentId,
	}

	msg := "failed to play audio"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return errors.WithMessage(err, msg)
	}

	if err := b.NotifyEventStream(payload, msg); err != nil {
		return errors.WithMessage(err, msg)
	}
	return nil
}

// Pause pauses audio playback.
func (c *Camera) Pause() error {
	payload := EventStreamPayload{
		Action:          "pause",
		Resource:        "audioPlayback/player",
		PublishResponse: false,
		From:            fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:              c.ParentId,
	}

	msg := "failed to pause audio"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return errors.WithMessage(err, msg)
	}

	if err := b.NotifyEventStream(payload, msg); err != nil {
		return errors.WithMessage(err, msg)
	}
	return nil
}

// Next moves audio playback to the next track.
func (c *Camera) Next() error {
	payload := EventStreamPayload{
		Action:          "nextTrack",
		Resource:        "audioPlayback/player",
		PublishResponse: false,
		From:            fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:              c.ParentId,
	}

	msg := "failed to skip audio"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return errors.WithMessage(err, msg)
	}

	if err := b.NotifyEventStream(payload, msg); err != nil {
		return errors.WithMessage(err, msg)
	}
	return nil
}

// Shuffle toggles the audio play back mode to shuffle or not.
func (c *Camera) Shuffle(on bool) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        "audioPlayback/config",
		PublishResponse: true,
		Properties: ShuffleProperties{
			Config: BaseShuffleProperties{
				ShuffleActive: on,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	var msg string
	if on {
		msg = "failed to enable shuffle"
	} else {
		msg = "failed to disable shuffle"
	}

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

func (c *Camera) Continuous() (response *EventStreamResponse, err error) {
	return c.SetLoopBackMode("continuous")
}

func (c *Camera) SingleTrack() (response *EventStreamResponse, err error) {
	return c.SetLoopBackMode("singleTrack")
}

func (c *Camera) SetLoopBackMode(loopbackMode string) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        "audioPlayback/config",
		PublishResponse: true,
		Properties: LoopbackModeProperties{
			Config: BaseLoopbackModeProperties{
				LoopbackMode: loopbackMode,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to set loop back mode to %s"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, fmt.Sprintf(msg, loopbackMode))
	}
	return b.makeEventStreamRequest(payload, msg)
}

func (c *Camera) GetAudioPlayback() (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "get",
		Resource:        "audioPlayback",
		PublishResponse: false,
		From:            fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:              c.ParentId,
	}

	msg := "failed to get audio playback"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

func (c *Camera) EnableSleepTimer(sleepTime int64 /* milliseconds */, sleepTimeRel int) (response *EventStreamResponse, err error) {
	if sleepTime == 0 {
		sleepTime = 300 + (time.Now().UnixNano() / 1000000) /* milliseconds */
	}

	payload := EventStreamPayload{
		Action:          "set",
		Resource:        "audioPlayback/config",
		PublishResponse: true,
		Properties: SleepTimerProperties{
			Config: BaseSleepTimerProperties{
				SleepTime:    sleepTime,
				SleepTimeRel: sleepTimeRel,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to enable sleep timer"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

func (c *Camera) DisableSleepTimer(sleepTimeRel int) (response *EventStreamResponse, err error) {
	if sleepTimeRel == 0 {
		sleepTimeRel = 300
	}

	payload := EventStreamPayload{
		Action:          "set",
		Resource:        "audioPlayback/config",
		PublishResponse: true,
		Properties: SleepTimerProperties{
			Config: BaseSleepTimerProperties{
				SleepTime:    0,
				SleepTimeRel: sleepTimeRel,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to disable sleep timer"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

/*
The follow methods are all related to the nightlight features of Arlo Baby.

NOTE: The current state is in: cameras[0]["properties"][0]["nightLight"] returned from the basestation.GetAssociatedCamerasState() method.
*/
func (c *Camera) NightLight(on bool) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: NightLightProperties{
			NightLight: BaseNightLightProperties{
				Enabled: on,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	var msg string
	if on {
		msg = "failed to turn night light on"
	} else {
		msg = "failed to turn night light off"
	}

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

func (c *Camera) SetNightLightBrightness(level int) (response *EventStreamResponse, err error) {
	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: NightLightProperties{
			NightLight: BaseNightLightProperties{
				Brightness: level,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to set night light brightness"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

// SetNightLightMode set the night light mode. Valid values are: "rainbow" or "rgb".
func (c *Camera) SetNightLightMode(mode string) (response *EventStreamResponse, err error) {
	msg := "failed to set night light brightness"

	if mode != "rainbow" && mode != "rgb" {
		return nil, errors.WithMessage(errors.New("mode can only be \"rainbow\" or \"rgb\""), msg)
	}

	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: NightLightProperties{
			NightLight: BaseNightLightProperties{
				Mode: mode,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

// SetNightLightColor sets the night light color to the RGB value specified by the three parameters, which have valid values from 0-255.
func (c *Camera) SetNightLightColor(red, blue, green int) (response *EventStreamResponse, err error) {
	// Sanity check; if the values are above or below the allowed limits, set them to their limit.
	if red < 0 {
		red = 0
	} else if red > 255 {
		red = 255
	}
	if blue < 0 {
		blue = 0
	} else if blue > 255 {
		blue = 255
	}
	if green < 0 {
		green = 0
	} else if green > 255 {
		green = 255
	}

	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: NightLightProperties{
			NightLight: BaseNightLightProperties{
				RGB: NightLightRGBProperties{
					Red:   red,
					Blue:  blue,
					Green: green,
				},
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to set night light color"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

func (c *Camera) EnableNightLightTimer(sleepTime int64 /* milliseconds */, sleepTimeRel int) (response *EventStreamResponse, err error) {
	if sleepTime == 0 {
		sleepTime = 300 + (time.Now().UnixNano() / 1000000) /* milliseconds */
	}

	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: NightLightProperties{
			NightLight: BaseNightLightProperties{
				SleepTime:    sleepTime,
				SleepTimeRel: sleepTimeRel,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to enable night light timer"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}

func (c *Camera) DisableNightLightTimer(sleepTimeRel int) (response *EventStreamResponse, err error) {
	if sleepTimeRel == 0 {
		sleepTimeRel = 300
	}

	payload := EventStreamPayload{
		Action:          "set",
		Resource:        fmt.Sprintf("cameras/%s", c.DeviceId),
		PublishResponse: true,
		Properties: NightLightProperties{
			NightLight: BaseNightLightProperties{
				SleepTime:    0,
				SleepTimeRel: sleepTimeRel,
			},
		},
		From: fmt.Sprintf("%s_%s", c.UserId, TransIdPrefix),
		To:   c.ParentId,
	}

	msg := "failed to disable night light timer"

	b := c.arlo.Basestations.Find(c.ParentId)
	if b == nil {
		err := fmt.Errorf("basestation (%s) not found for camera (%s)", c.ParentId, c.DeviceId)
		return nil, errors.WithMessage(err, msg)
	}
	return b.makeEventStreamRequest(payload, msg)
}
