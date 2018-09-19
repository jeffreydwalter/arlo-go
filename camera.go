package arlo

// A Camera is a Device of type "camera".
// This type is here just for semantics. Some methods explicitly require a device of a certain type.
type Camera Device

// Cameras is an array of Camera objects.
type Cameras []Camera
