package stlink

import (
	"errors"
	"fmt"

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

	mode, err := d.Mode()
	if err != nil {
		return err
	}

	if mode == StlinkModeDfu {
		err := d.ExitDFUMode()
		if err != nil {
			return nil
		}
	}

	mode, err = d.Mode()
	if err != nil {
		return err
	}
	if mode != StlinkModeDebug {
		err := d.EnterSWDMode()
		if err != nil {
			return err
		}
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

func (d *Device) String() string {
	s := ""
	name, err := d.Name()
	if err == nil {
		s += fmt.Sprintf(" name:    %s\n", name)
	} else {
		s += fmt.Sprintf(" name:    %v\n", err)
	}
	m, err := d.Mode()
	if err == nil {
		s += fmt.Sprintf(" mode:    %s\n", m)
	} else {
		s += fmt.Sprintf(" mode:    %v\n", err)
	}
	ver, err := d.Version()
	if err == nil {
		s += fmt.Sprintf(" version: %s\n", ver)
	} else {
		s += fmt.Sprintf(" version: %v\n", err)
	}
	v, err := d.TargetVoltage()
	if err == nil {
		s += fmt.Sprintf(" voltage: %.3f\n", v)
	} else {
		s += fmt.Sprintf(" voltage: %v\n", err)
	}
	if v < 1 {
		s += fmt.Sprintf(" target voltage too low!\n")
		return s
	}
	status, err := d.Status()
	if err == nil {
		s += fmt.Sprintf(" status:  %s\n", status)
	} else {
		s += fmt.Sprintf(" status:  %v\n", err)
	}
	cid, err := d.CoreID()
	if err == nil {
		s += fmt.Sprintf(" coreid:  %08x\n", cid)
	} else {
		s += fmt.Sprintf(" coreid:  %v\n", err)
	}
	cpu, err := d.CpuID()
	if err == nil {
		s += fmt.Sprintf(" cpu:     %08x\n", cpu)
	} else {
		s += fmt.Sprintf(" cpu:     %v\n", err)
	}
	chip, err := d.ChipID()
	if err == nil {
		s += fmt.Sprintf(" chip:    %08x\n", chip)
	} else {
		s += fmt.Sprintf(" chip:    %v\n", err)
	}
	dev, err := d.DevID()
	if err == nil {
		s += fmt.Sprintf(" dev:     %03x\n", dev)
	} else {
		s += fmt.Sprintf(" dev:     %v\n", err)
	}
	pn, err := d.CortexMPartNumber()
	if err == nil {
		s += fmt.Sprintf(" part-no: %s\n", pn)
	} else {
		s += fmt.Sprintf(" part-no: %v\n", err)
	}
	sz, err := d.FlashSize()
	if err == nil {
		s += fmt.Sprintf(" flash:   %d\n", sz)
	} else {
		s += fmt.Sprintf(" flash:   %v\n", err)
	}
	return s
}
