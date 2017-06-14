package stlink

import (
	"github.com/google/gousb"
)

const (
	stlinkUsbInEp     = 1
	stlinkUsbOutEpV2  = 2
	stlinkUsbOutEpV21 = 1
)

type Device struct {
	desc         *gousb.DeviceDesc
	dev          *gousb.Device
	interf       *gousb.Interface
	doneFunc     func()
	outEp        *gousb.OutEndpoint
	inEp         *gousb.InEndpoint
	SerialNumber string
	opened       bool
}

func (d *Device) init() error {
	var err error
	d.interf, d.doneFunc, err = d.dev.DefaultInterface()
	if err != nil {
		return err
	}
	d.inEp, err = d.interf.InEndpoint(stlinkUsbInEp)
	if err != nil {
		return err
	}

	ep := stlinkUsbOutEpV2
	if d.desc.Product == stlinkV21PID {
		ep = stlinkUsbOutEpV21
	}

	d.outEp, err = d.interf.OutEndpoint(ep)
	return err
}

func (d *Device) Close() error {
	if d.opened {
		if d.doneFunc != nil {
			d.doneFunc()
		}
		return d.dev.Close()
	}
	return nil
}

func (d *Device) write(b []byte) error {
	_, err := d.outEp.Write(b)
	return err
}

func (d *Device) read(n int) ([]byte, error) {
	rx := make([]byte, n, n)
	_, err := d.inEp.Read(rx)
	if err != nil {
		return nil, err
	}
	return rx, nil
}

const cmdSize = 16

type StlinkMode uint8

const (
	StlinkModeUnknown StlinkMode = 0xff
	StlinkModeDfu     StlinkMode = 0x00
	StlinkModeMass    StlinkMode = 0x01
	StlinkModeDebug   StlinkMode = 0x02
)

func (s StlinkMode) String() string {
	switch s {
	case StlinkModeDebug:
		return "debug"
	case StlinkModeDfu:
		return "dfu"
	case StlinkModeMass:
		return "mass"
	}
	return "unknown"
}

func (d *Device) GetMode() (StlinkMode, error) {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = stlinkCmdGetCurrentMode
	err := d.write(tx)
	if err != nil {
		return StlinkModeUnknown, err
	}
	rx, err := d.read(2)
	if err != nil {
		return StlinkModeUnknown, err
	}
	switch rx[0] {
	case uint8(StlinkModeDfu), uint8(StlinkModeMass), uint8(StlinkModeDebug):
		return StlinkMode(rx[0]), nil
	default:
		return StlinkModeUnknown, nil
	}
}

type StlinkStatus uint8

const (
	StlinkStatusUnknown     StlinkStatus = 0xff
	StlinkStatusCoreRunning              = 0x80
	StlinkStatusCoreHalted               = 0x81
)

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
