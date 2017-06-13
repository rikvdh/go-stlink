package stlink

import (
	"github.com/google/gousb"
)

type Device struct {
	desc         *gousb.DeviceDesc
	SerialNumber string
}
