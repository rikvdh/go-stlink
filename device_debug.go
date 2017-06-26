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

func (d *Device) Read16(addr uint32) (uint16, error) {
	if (addr % 2) != 0 {
		return 0, errors.New("unaligned access not allowed")
	}
	a, err := d.Read32(addr & 0xfffffffc)
	if err != nil {
		return 0, err
	}
	if addr%4 != 0 {
		a >>= 16
	}
	return uint16(a & 0xfffff), nil
}

const (
	cortexMIDCodeAddress   uint32 = 0xE0042000
	cortexMIDCodeM0Address uint32 = 0x40015800
)

func (d *Device) ChipID() (uint32, error) {
	pn, err := d.CortexMPartNumber()
	if err != nil {
		return 0, err
	}
	if pn == CortexMPartNumberM0 || pn == CortexMPartNumberM0Plus {
		return d.Read32(cortexMIDCodeM0Address)
	}
	return d.Read32(cortexMIDCodeAddress)
}

const cortexMCpuIDRegisterAddress uint32 = 0xE000ED00

func (d *Device) CpuID() (uint32, error) {
	id, err := d.Read32(cortexMCpuIDRegisterAddress)
	if err == nil {
		d.cpuID = id
	}
	return id, err
}

type ChipFamily uint16

const (
	ChipFamilySTM32F0             ChipFamily = 0x440
	ChipFamilySTM32F0Can          ChipFamily = 0x448
	ChipFamilySTM32F0Small        ChipFamily = 0x444
	ChipFamilySTM32F04            ChipFamily = 0x445
	ChipFamilySTM32F09X           ChipFamily = 0x442
	ChipFamilySTM32F1Connectivity ChipFamily = 0x418
	ChipFamilySTM32F1High         ChipFamily = 0x414
	ChipFamilySTM32F1Low          ChipFamily = 0x412
	ChipFamilySTM32F1Medium       ChipFamily = 0x410
	ChipFamilySTM32F1VLHigh       ChipFamily = 0x428
	ChipFamilySTM32F1VLMedium     ChipFamily = 0x420
	ChipFamilySTM32F1XL           ChipFamily = 0x430
	ChipFamilySTM32F2             ChipFamily = 0x411
	ChipFamilySTM32F3             ChipFamily = 0x422
	ChipFamilySTM32F3Small        ChipFamily = 0x439
	ChipFamilySTM32F303High       ChipFamily = 0x446
	ChipFamilySTM32F334           ChipFamily = 0x438
	ChipFamilySTM32F37x           ChipFamily = 0x432
	ChipFamilySTM32F4             ChipFamily = 0x413
	ChipFamilySTM32F4DE           ChipFamily = 0x433
	ChipFamilySTM32F4DSI          ChipFamily = 0x434
	ChipFamilySTM32F4HD           ChipFamily = 0x419
	ChipFamilySTM32F4LP           ChipFamily = 0x423
	ChipFamilySTM32F410           ChipFamily = 0x458
	ChipFamilySTM32F411RE         ChipFamily = 0x431
	ChipFamilySTM32F412           ChipFamily = 0x441
	ChipFamilySTM32F413           ChipFamily = 0x463
	ChipFamilySTM32F446           ChipFamily = 0x421
	ChipFamilySTM32F7             ChipFamily = 0x449
	ChipFamilySTM32F7Advanced     ChipFamily = 0x451
	ChipFamilySTM32F7Foundation   ChipFamily = 0x452
	ChipFamilySTM32L0             ChipFamily = 0x417
	ChipFamilySTM32L0Cat2         ChipFamily = 0x425
	ChipFamilySTM32L0Cat5         ChipFamily = 0x447
	ChipFamilySTM32L011           ChipFamily = 0x457
	ChipFamilySTM32L1Cat2         ChipFamily = 0x429
	ChipFamilySTM32L1High         ChipFamily = 0x436
	ChipFamilySTM32L1MediumLow    ChipFamily = 0x416
	ChipFamilySTM32L1MediumHigh   ChipFamily = 0x427
	ChipFamilySTM32L152RE         ChipFamily = 0x437
	ChipFamilySTM32L4             ChipFamily = 0x415
	ChipFamilySTM32L434X          ChipFamily = 0x435
	ChipFamilySTM32L4X6           ChipFamily = 0x461
)

func (d *Device) DevID() (ChipFamily, error) {
	id, err := d.ChipID()
	return ChipFamily(id & 0xfff), err
}

func (d *Device) FlashSize() (uint16, error) {
	pn, err := d.DevID()
	if err != nil {
		return 0, err
	}
	switch pn {
	case ChipFamilySTM32F0, ChipFamilySTM32F09X, ChipFamilySTM32F0Can,
		ChipFamilySTM32F0Small, ChipFamilySTM32F04, ChipFamilySTM32F3,
		ChipFamilySTM32F37x, ChipFamilySTM32F334, ChipFamilySTM32F3Small,
		ChipFamilySTM32F303High:
		return d.Read16(0x1ffff7cc)

	case ChipFamilySTM32F1Medium, ChipFamilySTM32F1High, ChipFamilySTM32F1Low,
		ChipFamilySTM32F1Connectivity, ChipFamilySTM32F1VLMedium, ChipFamilySTM32F1VLHigh,
		ChipFamilySTM32F1XL:
		return d.Read16(0x1ffff7e0)

	case ChipFamilySTM32F2, ChipFamilySTM32F4, ChipFamilySTM32F4HD,
		ChipFamilySTM32F446, ChipFamilySTM32F411RE, ChipFamilySTM32F4DE,
		ChipFamilySTM32F4DSI, ChipFamilySTM32F412, ChipFamilySTM32F410,
		ChipFamilySTM32F413, ChipFamilySTM32F7Foundation, ChipFamilySTM32F4LP:
		return d.Read16(0x1fff7a22)

	case ChipFamilySTM32L4, ChipFamilySTM32L434X, ChipFamilySTM32L4X6:
		return d.Read16(0x1fff75e0)

	case ChipFamilySTM32L011, ChipFamilySTM32L0Cat2, ChipFamilySTM32L0,
		ChipFamilySTM32L0Cat5:
		return d.Read16(0x1ff8007c)

	case ChipFamilySTM32L1MediumLow, ChipFamilySTM32L1Cat2:
		return d.Read16(0x1ff8004c)

	case ChipFamilySTM32L1MediumHigh, ChipFamilySTM32L1High, ChipFamilySTM32L152RE:
		return d.Read16(0x1ff800cc)

	case ChipFamilySTM32F7, ChipFamilySTM32F7Advanced:
		return d.Read16(0x1ff0f442)

	}
	return 0, errors.New("unknown core")
}
