package stlink

import (
	"encoding/binary"
	"errors"
)

type StlinkStatus uint8

const (
	StlinkStatusUnknown     StlinkStatus = 0xff
	StlinkStatusCoreRunning StlinkStatus = 0x80
	StlinkStatusCoreHalted  StlinkStatus = 0x81
)

func (s StlinkStatus) String() string {
	switch s {
	case StlinkStatusCoreHalted:
		return "halted"
	case StlinkStatusCoreRunning:
		return "running"
	}
	return "unknown"
}

func (d *Device) GetStatus() (StlinkStatus, error) {
	d.coreState = StlinkStatusUnknown
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdDebug)
	tx[1] = byte(stlinkCmdDebugGetStatus)
	err := d.write(tx)
	if err != nil {
		return StlinkStatusUnknown, err
	}

	rx, err := d.read(2)
	if err != nil {
		return StlinkStatusUnknown, err
	}
	if rx[0] == byte(StlinkStatusCoreRunning) || rx[0] == byte(StlinkStatusCoreHalted) {
		d.coreState = StlinkStatus(rx[0])
		return d.coreState, nil
	}
	return StlinkStatusUnknown, nil
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

func (d *Device) GetCoreID() (uint32, error) {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdDebug)
	tx[1] = byte(stlinkCmdDebugReadCoreid)
	err := d.write(tx)
	if err != nil {
		return 0, err
	}

	var v uint32
	if err := binary.Read(d.inEp, binary.LittleEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func (d *Device) Halt() error {
	return d.Step()
}

func (d *Device) Step() error {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdDebug)
	tx[1] = byte(stlinkCmdDebugStepCore)
	err := d.write(tx)
	if err != nil {
		return err
	}
	_, err = d.read(2)
	return err
}

func (d *Device) ForceDebug() error {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdDebug)
	tx[1] = byte(stlinkCmdDebugForce)
	err := d.write(tx)
	if err != nil {
		return err
	}
	_, err = d.read(2)
	return err
}

func (d *Device) Reset() error {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdDebug)
	tx[1] = byte(stlinkCmdDebugResetsys)
	err := d.write(tx)
	if err != nil {
		return err
	}
	_, err = d.read(2)
	return err
}

func (d *Device) HardReset() error {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdDebug)
	tx[1] = byte(stlinkCmdDebugHardReset)
	err := d.write(tx)
	if err != nil {
		return err
	}
	_, err = d.read(2)
	return err
}

func (d *Device) Run() error {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdDebug)
	tx[1] = byte(stlinkCmdDebugRunCore)
	err := d.write(tx)
	if err != nil {
		return err
	}
	_, err = d.read(2)
	return err
}
