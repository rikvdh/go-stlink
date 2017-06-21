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

func (d *Device) Status() (StlinkStatus, error) {
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

func (d *Device) ClockSpeed() (StlinkClockSpeed, error) {
	return StlinkClockSpeed5, errors.New("not implemented")
}

func (d *Device) CoreID() (uint32, error) {
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

func (d *Device) Write32(addr, w uint32) error {
	if d.coreState == StlinkStatusCoreRunning {
		d.Halt()
	}
	defer func() {
		if d.coreState == StlinkStatusCoreRunning {
			d.Run()
		}
	}()
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdDebug)
	tx[1] = byte(stlinkCmdDebugJtagWritedebug32bit)

	binary.LittleEndian.PutUint32(tx[2:], addr)
	binary.LittleEndian.PutUint32(tx[6:], w)

	err := d.write(tx)
	if err != nil {
		return err
	}
	_, err = d.read(8)
	return err
}

func (d *Device) Read32(addr uint32) (uint32, error) {
	tx := make([]byte, 6, 6)
	tx[0] = byte(stlinkCmdDebug)
	tx[1] = byte(stlinkCmdDebugJtagReaddebug32bit)

	binary.LittleEndian.PutUint32(tx[2:], addr)

	err := d.write(tx)
	if err != nil {
		return 0, err
	}
	var v [2]uint32
	if err := binary.Read(d.inEp, binary.LittleEndian, &v); err != nil {
		return 0, err
	}
	return v[1], nil
}

const (
	cortexMIDCodeAddress   uint32 = 0xE0042000
	cortexMIDCodeM0Address uint32 = 0x40015800
)

func (d *Device) ChipID() (uint32, error) {
	return 0, nil
	//stlink2_get_cpuid(dev);
	//enum stlink2_cortexm_cpuid_partno partno = stlink2_cortexm_cpuid_get_partno(dev->mcu.cpuid);

	/*if (partno == STLINK2_CORTEXM_CPUID_PARTNO_M0 ||
	    partno == STLINK2_CORTEXM_CPUID_PARTNO_M0_PLUS)
		stlink2_read_debug32(dev, cortexMIDCodeM0Address, &dev->mcu.chipid);
	else*/
	//	stlink2_read_debug32(dev, cortexMIDCodeAddress, &dev->mcu.chipid);

	//return dev->mcu.chipid;
}

const cortexMCpuIDRegisterAddress uint32 = 0xE000ED00

func (d *Device) CpuID() (uint32, error) {
	return d.Read32(cortexMCpuIDRegisterAddress)
}

func (d *Device) DevID() (uint16, error) {
	id, err := d.CpuID()
	return uint16(id & 0xfff), err
}
