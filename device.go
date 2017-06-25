package stlink

import (
	"errors"

	"github.com/google/gousb"
)

const (
	cmdSize           = 16
	stlinkUsbInEp     = 1
	stlinkUsbOutEpV2  = 2
	stlinkUsbOutEpV21 = 1
)

// Device represents a ST-link device
type Device struct {
	desc         *gousb.DeviceDesc
	dev          *gousb.Device
	interf       *gousb.Interface
	doneFunc     func()
	outEp        *gousb.OutEndpoint
	inEp         *gousb.InEndpoint
	SerialNumber string
	opened       bool
	coreState    StlinkStatus
	cpuID        uint32
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
	if err != nil {
		return err
	}
	_, err = d.Status()
	return err
}

// Close closes the device when needed
func (d *Device) Close() error {
	if d.opened {
		if d.doneFunc != nil {
			d.doneFunc()
		}
		return d.dev.Close()
	}
	return nil
}

func (d *Device) Name() (string, error) {
	switch d.desc.Product {
	case stlinkV21PID:
		return "ST-link V2-1", nil
	case stlinkV1PID:
		return "ST-link V1", nil
	case stlinkV2PID:
		return "ST-link V2", nil
	}
	return "", errors.New("unknown device")
}

func (d *Device) write(b []byte) error {
	if !d.opened {
		return errors.New("device closed")
	}
	_, err := d.outEp.Write(b)
	return err
}

func (d *Device) read(n int) ([]byte, error) {
	if !d.opened {
		return nil, errors.New("device closed")
	}
	rx := make([]byte, n, n)
	_, err := d.inEp.Read(rx)
	if err != nil {
		return nil, err
	}
	return rx, nil
}
