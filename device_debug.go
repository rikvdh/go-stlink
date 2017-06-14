package stlink

import "errors"

type StlinkStatus uint8

const (
	StlinkStatusUnknown     StlinkStatus = 0xff
	StlinkStatusCoreRunning              = 0x80
	StlinkStatusCoreHalted               = 0x81
)

func (d *Device) GetStatus() (StlinkStatus, error) {
	return StlinkStatusUnknown, errors.New("not implemented")
}

type StlinkClockSpeed uint16

const (
	StlinkClockSpeed4000 StlinkClockSpeed = 0
	StlinkClockSpeed1800                  = 1
	StlinkClockSpeed1200                  = 2
	StlinkClockSpeed950                   = 3
	StlinkClockSpeed480                   = 7
	StlinkClockSpeed240                   = 15
	StlinkClockSpeed125                   = 31
	StlinkClockSpeed100                   = 40
	StlinkClockSpeed50                    = 79
	StlinkClockSpeed25                    = 158
	StlinkClockSpeed15                    = 265
	StlinkClockSpeed5                     = 798
)

func (d *Device) GetClockSpeed() (StlinkClockSpeed, error) {
	return StlinkClockSpeed5, errors.New("not implemented")
}
