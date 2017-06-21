package stlink

import (
	"encoding/binary"
	"errors"
	"fmt"
)

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

// GetMode reads the mode for an ST-link, see StlinkMode
func (d *Device) Mode() (StlinkMode, error) {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdGetCurrentMode)
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

func (d *Device) Version() (string, error) {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdGetVersion)
	err := d.write(tx)
	if err != nil {
		return "", err
	}
	rx, err := d.read(6)
	if err != nil {
		return "", err
	}

	stlink := uint8((rx[0] & 0xf0) >> 4)
	jtag := uint8(((rx[0] & 0x0f) << 2) | ((rx[1] & 0xc0) >> 6))
	swim := uint8(rx[1] & 0x3f)

	return fmt.Sprintf("V%dJ%dS%d", stlink, jtag, swim), nil
}

func (d *Device) TargetVoltage() (float32, error) {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdGetTargetVoltage)
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
