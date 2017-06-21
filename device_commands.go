package stlink

import (
	"encoding/binary"
	"errors"
)

type StlinkMode uint8

const (
	StlinkModeUnknown StlinkMode = 0xff
	StlinkModeDfu                = 0x00
	StlinkModeMass               = 0x01
	StlinkModeDebug              = 0x02
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

// GetMode reads the mode for an ST-link, see StlinkMode
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

func (d *Device) GetTargetVoltage() (float32, error) {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = stlinkCmdGetTargetVoltage
	err := d.write(tx)
	if err != nil {
		return 0, err
	}

	var v [2]int32
	if err := binary.Read(d.inEp, binary.LittleEndian, &v); err != nil {
		return 0, err
	}

	if v[0] == 0 {
		return 0, errors.New("measured voltage is zero")
	}
	return 2.4 * float32(v[1]) / float32(v[0]), nil
}
